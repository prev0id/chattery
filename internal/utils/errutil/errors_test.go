package errutil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKind_StatusCode(t *testing.T) {
	t.Parallel()

	type expected struct {
		value int
	}
	tests := []struct {
		name     string
		expected expected
		k        Kind
	}{
		{
			name: "internal",
			k:    Internal,
			expected: expected{
				value: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid_request",
			k:    InvalidRequest,
			expected: expected{
				value: http.StatusBadRequest,
			},
		},
		{
			name: "unauthorized",
			k:    Unauthorized,
			expected: expected{
				value: http.StatusUnauthorized,
			},
		},
		{
			name: "permission",
			k:    Permission,
			expected: expected{
				value: http.StatusForbidden,
			},
		},
		{
			name: "exist",
			k:    Exist,
			expected: expected{
				value: http.StatusConflict,
			},
		},
		{
			name: "not_found",
			k:    NotFound,
			expected: expected{
				value: http.StatusNotFound,
			},
		},
		{
			name: "unknown",
			k:    Kind(100),
			expected: expected{
				value: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.k.StatusCode()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestKind_String(t *testing.T) {
	t.Parallel()

	type expected struct {
		value string
	}
	tests := []struct {
		name     string
		expected expected
		k        Kind
	}{
		{
			name: "internal",
			k:    Internal,
			expected: expected{
				value: "internal error",
			},
		},
		{
			name: "invalid_request",
			k:    InvalidRequest,
			expected: expected{
				value: "invalid request",
			},
		},
		{
			name: "unauthorized",
			k:    Unauthorized,
			expected: expected{
				value: "unauthorized",
			},
		},
		{
			name: "permission",
			k:    Permission,
			expected: expected{
				value: "forbidden",
			},
		},
		{
			name: "exist",
			k:    Exist,
			expected: expected{
				value: "already exists",
			},
		},
		{
			name: "not_found",
			k:    NotFound,
			expected: expected{
				value: "not found",
			},
		},
		{
			name: "unknown",
			k:    Kind(100),
			expected: expected{
				value: "unknown error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.k.String()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestE(t *testing.T) {
	t.Parallel()

	existingErr := &Error{kind: NotFound, message: "missing"}

	type args struct {
		errs []error
	}
	type expected struct {
		value *Error
	}
	tests := []struct {
		expected expected
		name     string
		args     args
	}{
		{
			name: "empty",
			expected: expected{
				value: &Error{},
			},
		},
		{
			name: "wraps_regular_error",
			args: args{
				errs: []error{assert.AnError},
			},
			expected: expected{
				value: &Error{err: assert.AnError},
			},
		},
		{
			name: "returns_existing_error",
			args: args{
				errs: []error{existingErr},
			},
			expected: expected{
				value: existingErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := E(tt.args.errs...)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type expected struct {
		value string
	}
	tests := []struct {
		expected expected
		name     string
		fields   fields
	}{
		{
			name: "kind_only",
			fields: fields{
				kind: Internal,
			},
			expected: expected{
				value: "internal error",
			},
		},
		{
			name: "kind_and_message",
			fields: fields{
				kind:    NotFound,
				message: "user not found",
			},
			expected: expected{
				value: "not found: user not found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.Error()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_Kind(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type args struct {
		kind Kind
	}
	type expected struct {
		value *Error
	}
	tests := []struct {
		expected expected
		name     string
		fields   fields
		args     args
	}{
		{
			name: "sets_kind",
			args: args{
				kind: Permission,
			},
			expected: expected{
				value: &Error{kind: Permission},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.Kind(tt.args.kind)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_GetKind(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type expected struct {
		value Kind
	}
	tests := []struct {
		name     string
		fields   fields
		expected expected
	}{
		{
			name: "returns_kind",
			fields: fields{
				kind: NotFound,
			},
			expected: expected{
				value: NotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.GetKind()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_Debug(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type args struct {
		messages []string
	}
	type expected struct {
		value *Error
	}
	tests := []struct {
		expected expected
		name     string
		args     args
		fields   fields
	}{
		{
			name: "appends_debug",
			fields: fields{
				debug: []string{"first"},
			},
			args: args{
				messages: []string{"second"},
			},
			expected: expected{
				value: &Error{debug: []string{"first", "second"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.Debug(tt.args.messages...)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_GetDebug(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type expected struct {
		value []string
	}
	tests := []struct {
		name     string
		expected expected
		fields   fields
	}{
		{
			name: "returns_debug",
			fields: fields{
				debug: []string{"first", "second"},
			},
			expected: expected{
				value: []string{"first", "second"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.GetDebug()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_Message(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type args struct {
		message string
	}
	type expected struct {
		value *Error
	}
	tests := []struct {
		expected expected
		args     args
		name     string
		fields   fields
	}{
		{
			name: "sets_message",
			args: args{
				message: "user not found",
			},
			expected: expected{
				value: &Error{message: "user not found"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.Message(tt.args.message)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_Messagef(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type args struct {
		format string
		args   []any
	}
	type expected struct {
		value *Error
	}
	tests := []struct {
		expected expected
		args     args
		name     string
		fields   fields
	}{
		{
			name: "formats_message",
			args: args{
				format: "user %s not found",
				args:   []any{"alice"},
			},
			expected: expected{
				value: &Error{message: "user alice not found"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.Messagef(tt.args.format, tt.args.args...)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_GetMessage(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type expected struct {
		value string
	}
	tests := []struct {
		expected expected
		name     string
		fields   fields
	}{
		{
			name: "returns_message",
			fields: fields{
				message: "user not found",
			},
			expected: expected{
				value: "user not found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.GetMessage()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestError_GetError(t *testing.T) {
	t.Parallel()

	type fields struct {
		err     error
		message string
		debug   []string
		kind    Kind
	}
	type expected struct {
		value error
	}
	tests := []struct {
		expected expected
		name     string
		fields   fields
	}{
		{
			name: "returns_error",
			fields: fields{
				err: assert.AnError,
			},
			expected: expected{
				value: assert.AnError,
			},
		},
		{
			name: "returns_nil",
			expected: expected{
				value: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &Error{
				err:     tt.fields.err,
				message: tt.fields.message,
				debug:   tt.fields.debug,
				kind:    tt.fields.kind,
			}

			value := e.GetError()

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestIs(t *testing.T) {
	t.Parallel()

	type args struct {
		err  error
		kind Kind
	}
	type expected struct {
		value bool
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "same_kind",
			args: args{
				kind: NotFound,
				err:  E().Kind(NotFound),
			},
			expected: expected{
				value: true,
			},
		},
		{
			name: "different_kind",
			args: args{
				kind: InvalidRequest,
				err:  E().Kind(NotFound),
			},
			expected: expected{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := Is(tt.args.kind, tt.args.err)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}
