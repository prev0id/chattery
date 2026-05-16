package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	user_api "chattery/internal/api/user"
	"chattery/internal/domain"
)

const baseURL = "http://localhost:8080"

const requestTimeout = 10 * time.Second

type Response struct {
	Header     http.Header
	Cookies    []*http.Cookie
	Body       []byte
	StatusCode int
}

func MustPostCreateUser(t testing.TB, request *user_api.PostCreateUserRequest, cookies ...*http.Cookie) *Response {
	t.Helper()

	return mustDoJSON(t, http.MethodPost, "/v1/user/create", request, cookies...)
}

func MustPostLogin(t testing.TB, request *user_api.PostLoginRequest, cookies ...*http.Cookie) *Response {
	t.Helper()

	return mustDoJSON(t, http.MethodPost, "/v1/user/login", request, cookies...)
}

func MustPostLogout(t testing.TB, cookies ...*http.Cookie) *Response {
	t.Helper()

	return mustDo(t, http.MethodPost, "/v1/user/logout", nil, "", cookies...)
}

func MustDeleteMe(t testing.TB, cookies ...*http.Cookie) *Response {
	t.Helper()

	return mustDo(t, http.MethodDelete, "/v1/user/delete", nil, "", cookies...)
}

func MustGetMe(t testing.TB, cookies ...*http.Cookie) *Response {
	t.Helper()

	return mustDo(t, http.MethodGet, "/v1/user/me", nil, "", cookies...)
}

func (r *Response) RequireStatus(t testing.TB, status int) {
	t.Helper()

	require.Equalf(t, status, r.StatusCode, "response body: %s", string(r.Body))
}

func (r *Response) RequireJSON(t testing.TB, value any) {
	t.Helper()

	require.NoError(t, json.Unmarshal(r.Body, value))
}

func (r *Response) RequireSessionCookie(t testing.TB) *http.Cookie {
	t.Helper()

	cookie := r.sessionCookie(t)
	require.NotEmpty(t, cookie.Value)

	return cookie
}

func (r *Response) RequireNoSessionCookie(t testing.TB) {
	t.Helper()

	for _, cookie := range r.Cookies {
		require.NotEqual(t, domain.SessionCookieName, cookie.Name)
	}
}

func (r *Response) RequireClearedSessionCookie(t testing.TB) {
	t.Helper()

	cookie := r.sessionCookie(t)
	require.Empty(t, cookie.Value)
	require.True(t, cookie.Expires.Before(time.Now()))
}

func (r *Response) RequireErrorContains(t testing.TB, expected string) {
	t.Helper()

	var body struct {
		Message string `json:"message"`
	}
	r.RequireJSON(t, &body)
	require.Contains(t, body.Message, expected)
}

func (r *Response) sessionCookie(t testing.TB) *http.Cookie {
	t.Helper()

	for _, cookie := range r.Cookies {
		if cookie.Name == domain.SessionCookieName {
			return cookie
		}
	}

	require.FailNowf(t, "session cookie not found", "cookies: %v", r.Cookies)
	return nil
}

func mustDoJSON(t testing.TB, method, path string, request any, cookies ...*http.Cookie) *Response {
	t.Helper()

	raw, err := json.Marshal(request)
	require.NoError(t, err)

	return mustDo(t, method, path, bytes.NewReader(raw), "application/json", cookies...)
}

func mustDo(t testing.TB, method, path string, body io.Reader, contentType string, cookies ...*http.Cookie) *Response {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	require.NoError(t, err)

	request.Header.Set("Accept", "application/json")
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	for _, cookie := range cookies {
		if cookie != nil {
			request.AddCookie(cookie)
		}
	}

	response, err := httpClient().Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return &Response{
		StatusCode: response.StatusCode,
		Header:     response.Header.Clone(),
		Cookies:    response.Cookies(),
		Body:       raw,
	}
}

func httpClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
