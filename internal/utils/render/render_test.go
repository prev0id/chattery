package render

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

type renderTestValue struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func TestTimestamp(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 15, 4, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)

	type args struct {
		t time.Time
	}
	type expected struct {
		value string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "today",
			args: args{
				t: today,
			},
			expected: expected{
				value: "Today, 15:04",
			},
		},
		{
			name: "yesterday",
			args: args{
				t: yesterday,
			},
			expected: expected{
				value: "Yesterday, 15:04",
			},
		},
		{
			name: "older",
			args: args{
				t: time.Date(2000, time.January, 2, 3, 4, 0, 0, time.UTC),
			},
			expected: expected{
				value: "Jan 2, 03:04",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := Timestamp(tt.args.t)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		value any
	}
	type expected struct {
		body        string
		contentType string
		statusCode  int
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				value: renderTestValue{Name: "alice", ID: 7},
			},
			expected: expected{
				body:        `{"name":"alice","id":7}`,
				contentType: "application/json",
				statusCode:  http.StatusOK,
			},
		},
		{
			name: "marshal_error",
			args: args{
				value: func() {},
			},
			expected: expected{
				body:        `{"message":"internal error"}`,
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			JSON(recorder, request, tt.args.value)

			assert.Equal(t, tt.expected.statusCode, recorder.Code)
			assert.Equal(t, tt.expected.contentType, recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.expected.body, recorder.Body.String())
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}
	type expected struct {
		body        string
		contentType string
		statusCode  int
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "not_found",
			args: args{
				err: errutil.E().Kind(errutil.NotFound).Message("user not found"),
			},
			expected: expected{
				body:        `{"message":"not found: user not found"}`,
				contentType: "application/json",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name: "internal",
			args: args{
				err: assert.AnError,
			},
			expected: expected{
				body:        `{"message":"assert.AnError general error for testing"}`,
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			Error(recorder, request, tt.args.err)

			assert.Equal(t, tt.expected.statusCode, recorder.Code)
			assert.Equal(t, tt.expected.contentType, recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.expected.body, recorder.Body.String())
		})
	}
}

func Test_setContentTypeJSON(t *testing.T) {
	t.Parallel()

	type expected struct {
		contentType string
	}
	tests := []struct {
		name     string
		expected expected
	}{
		{
			name: "sets_content_type",
			expected: expected{
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			setContentTypeJSON(recorder)

			assert.Equal(t, tt.expected.contentType, recorder.Header().Get("Content-Type"))
		})
	}
}

func TestJSONBytes(t *testing.T) {
	t.Parallel()

	type args struct {
		value any
	}
	type expected struct {
		value    []byte
		errDebug []string
		err      bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				value: renderTestValue{Name: "alice", ID: 7},
			},
			expected: expected{
				value: []byte(`{"name":"alice","id":7}`),
			},
		},
		{
			name: "marshal_error",
			args: args{
				value: func() {},
			},
			expected: expected{
				err:      true,
				errDebug: []string{"json.Marshal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := JSONBytes(tt.args.value)

			requireRenderError(t, err, tt.expected.err, tt.expected.errDebug)
			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestJSONString(t *testing.T) {
	t.Parallel()

	type args struct {
		value any
	}
	type expected struct {
		value    string
		errDebug []string
		err      bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				value: renderTestValue{Name: "alice", ID: 7},
			},
			expected: expected{
				value: `{"name":"alice","id":7}`,
			},
		},
		{
			name: "marshal_error",
			args: args{
				value: func() {},
			},
			expected: expected{
				err:      true,
				errDebug: []string{"json.Marshal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := JSONString(tt.args.value)

			requireRenderError(t, err, tt.expected.err, tt.expected.errDebug)
			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestEvent(t *testing.T) {
	t.Parallel()

	type args struct {
		payload   any
		eventType event_desc.Type
		channel   event_desc.Channel
	}
	type expected struct {
		body     string
		errDebug []string
		err      bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				eventType: event_desc.TypeMessage,
				channel: event_desc.Channel{
					Type: event_desc.ChannelDM,
					ID:   42,
				},
				payload: renderTestValue{Name: "alice", ID: 7},
			},
			expected: expected{
				body: `{"type":"message","channel":{"type":"dm","id":42},"payload":{"name":"alice","id":7}}`,
			},
		},
		{
			name: "payload_marshal_error",
			args: args{
				eventType: event_desc.TypeMessage,
				payload:   func() {},
			},
			expected: expected{
				err:      true,
				errDebug: []string{"json.Marshal", "JSONBytes", "payload"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := Event(tt.args.eventType, tt.args.channel, tt.args.payload)

			requireRenderError(t, err, tt.expected.err, tt.expected.errDebug)
			if tt.expected.body != "" {
				assert.JSONEq(t, tt.expected.body, string(value))
			} else {
				assert.Equal(t, []byte(nil), value)
			}
		})
	}
}

func TestAvatarURL(t *testing.T) {
	t.Parallel()

	type args struct {
		username domain.Username
	}
	type expected struct {
		value string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "username",
			args: args{
				username: "alice",
			},
			expected: expected{
				value: "/v1/image/alice.jpeg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := AvatarURL(tt.args.username)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func requireRenderError(t *testing.T, err error, expectedErr bool, expectedDebug []string) {
	t.Helper()

	if !expectedErr {
		require.NoError(t, err)
		return
	}

	require.Error(t, err)
	assert.Equal(t, expectedDebug, errutil.E(err).GetDebug())
}
