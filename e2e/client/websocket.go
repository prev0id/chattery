package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"

	"chattery/internal/api/websocket/event_desc"
)

const baseWSURL = "ws://localhost:8080"

func MustDialWebsocket(t testing.TB, cookies ...*http.Cookie) *websocket.Conn {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	conn, response, err := websocket.Dial(ctx, baseWSURL+"/ws/", websocketDialOptions(cookies...))
	require.NoErrorf(t, err, "response: %+v", response)
	require.NotNil(t, conn)

	return conn
}

func MustRejectWebsocket(t testing.TB, expectedStatus int, cookies ...*http.Cookie) *Response {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	conn, response, err := websocket.Dial(ctx, baseWSURL+"/ws/", websocketDialOptions(cookies...))
	require.Error(t, err)
	require.Nil(t, conn)
	require.NotNil(t, response)
	defer response.Body.Close()
	require.Equal(t, expectedStatus, response.StatusCode)

	raw, readErr := io.ReadAll(response.Body)
	require.NoError(t, readErr)

	return &Response{
		StatusCode: response.StatusCode,
		Header:     response.Header.Clone(),
		Cookies:    response.Cookies(),
		Body:       raw,
	}
}

func MustWriteWSEvent(t testing.TB, conn *websocket.Conn, event *event_desc.Event) {
	t.Helper()

	raw, err := json.Marshal(event)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	require.NoError(t, conn.Write(ctx, websocket.MessageText, raw))
}

func MustReadWSEvent(t testing.TB, conn *websocket.Conn, timeout time.Duration) *event_desc.Event {
	t.Helper()

	event, err := readWSEvent(conn, timeout)
	require.NoError(t, err)

	return event
}

func RequireNoWSEvent(t testing.TB, conn *websocket.Conn, timeout time.Duration) {
	t.Helper()

	event, err := readWSEvent(conn, timeout)
	if errors.Is(err, context.DeadlineExceeded) {
		return
	}

	require.NoError(t, err)
	require.Failf(t, "unexpected websocket event", "event=%+v", event)
}

func DrainWSEvents(t testing.TB, conn *websocket.Conn, timeout time.Duration) {
	t.Helper()

	for {
		_, err := readWSEvent(conn, timeout)
		if errors.Is(err, context.DeadlineExceeded) {
			return
		}
		require.NoError(t, err)
	}
}

func MustCloseWebsocket(t testing.TB, conn *websocket.Conn) {
	t.Helper()

	if conn == nil {
		return
	}

	require.NoError(t, conn.Close(websocket.StatusNormalClosure, ""))
}

func websocketDialOptions(cookies ...*http.Cookie) *websocket.DialOptions {
	header := http.Header{}
	request := http.Request{Header: header}

	for _, cookie := range cookies {
		if cookie != nil {
			request.AddCookie(cookie)
		}
	}

	return &websocket.DialOptions{
		HTTPHeader: header,
	}
}

func readWSEvent(conn *websocket.Conn, timeout time.Duration) (*event_desc.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	messageType, raw, err := conn.Read(ctx)
	if err != nil {
		return nil, err
	}
	if messageType != websocket.MessageText {
		return nil, errors.New("unexpected websocket message type")
	}

	var event event_desc.Event
	if err := json.Unmarshal(raw, &event); err != nil {
		return nil, err
	}

	return &event, nil
}
