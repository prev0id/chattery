package logger

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

var testLogStore = newCaptureHandler()

func init() {
	slog.SetDefault(slog.New(testLogStore))
}

func TestError(t *testing.T) {
	t.Parallel()

	type args struct {
		err     error
		message string
		attr    []slog.Attr
	}
	type expected struct {
		message      string
		scope        string
		errKind      string
		errMessage   string
		errText      string
		errDebug     []string
		recordExists bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "records_error",
			args: args{
				err: errutil.E(assert.AnError).
					Kind(errutil.NotFound).
					Message("user not found").
					Debug("adapter.UserByID"),
				message: "test_logger_error_records_error",
				attr: []slog.Attr{
					slog.String("scope", "test"),
				},
			},
			expected: expected{
				message:      "test_logger_error_records_error",
				scope:        "test",
				errKind:      "not found",
				errMessage:   "user not found",
				errDebug:     []string{"adapter.UserByID"},
				errText:      assert.AnError.Error(),
				recordExists: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			Error(tt.args.err, tt.args.message, tt.args.attr...)
			record, ok := testLogStore.record(tt.expected.message)

			if !assert.Equal(t, tt.expected.recordExists, ok) {
				return
			}
			assert.Equal(t, tt.expected.scope, record.attrs["scope"].String())
			errorAttrs := record.attrs["error"].Group()
			assert.Equal(t, tt.expected.errKind, attrValue(errorAttrs, "kind").String())
			assert.Equal(t, tt.expected.errMessage, attrValue(errorAttrs, "message").String())
			assert.Equal(t, tt.expected.errText, fmt.Sprint(attrValue(errorAttrs, "err").Any()))
			assert.Equal(t, tt.expected.errDebug, attrValue(errorAttrs, "debug").Any())
		})
	}
}

func TestErrorCtx(t *testing.T) {
	t.Parallel()

	ctx := domain.UserIDToContext(context.Background(), 42)
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "request-id")

	type args struct {
		ctx     context.Context
		err     error
		message string
		attr    []slog.Attr
	}
	type expected struct {
		requestID    string
		message      string
		userID       int64
		recordExists bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "records_context_fields",
			args: args{
				ctx:     ctx,
				err:     errutil.E(assert.AnError),
				message: "test_logger_error_ctx_records_context_fields",
			},
			expected: expected{
				requestID:    "request-id",
				message:      "test_logger_error_ctx_records_context_fields",
				userID:       42,
				recordExists: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ErrorCtx(tt.args.ctx, tt.args.err, tt.args.message, tt.args.attr...)
			record, ok := testLogStore.record(tt.expected.message)

			if !assert.Equal(t, tt.expected.recordExists, ok) {
				return
			}
			assert.Equal(t, tt.expected.requestID, record.attrs["request_id"].String())
			assert.Equal(t, tt.expected.userID, record.attrs["user_id"].Int64())
		})
	}
}

type capturedRecord struct {
	attrs   map[string]slog.Value
	message string
}

type captureHandler struct {
	records map[string]capturedRecord
	attrs   []slog.Attr
	m       sync.Mutex
}

func newCaptureHandler() *captureHandler {
	return &captureHandler{
		records: make(map[string]capturedRecord),
	}
}

func (*captureHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *captureHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]slog.Value)
	for _, attr := range h.attrs {
		attrs[attr.Key] = attr.Value
	}
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value
		return true
	})

	h.m.Lock()
	defer h.m.Unlock()
	h.records[record.Message] = capturedRecord{
		attrs:   attrs,
		message: record.Message,
	}
	return nil
}

func (h *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := &captureHandler{
		records: h.records,
		attrs:   append(h.attrs, attrs...),
	}
	return next
}

func (h *captureHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *captureHandler) record(message string) (capturedRecord, bool) {
	h.m.Lock()
	defer h.m.Unlock()
	record, ok := h.records[message]
	return record, ok
}

func attrValue(attrs []slog.Attr, key string) slog.Value {
	for _, attr := range attrs {
		if attr.Key == key {
			return attr.Value
		}
	}
	return slog.Value{}
}
