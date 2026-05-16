package bind

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chattery/internal/utils/errutil"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	validValue := &jsonTestValue{Name: "alice", ID: 7}

	type args struct {
		request *http.Request
	}
	type expected struct {
		value    *jsonTestValue
		errDebug []string
		errKind  errutil.Kind
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"alice","id":7}`)),
			},
			expected: expected{
				value: validValue,
			},
		},
		{
			name: "body_read_error",
			args: args{
				request: requestWithBody(errorReadCloser{err: assert.AnError}),
			},
			expected: expected{
				errKind:  errutil.InvalidRequest,
				errDebug: []string{"io.ReadAll"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := JSON[jsonTestValue](tt.args.request)

			if tt.expected.errDebug == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.True(t, errutil.Is(tt.expected.errKind, err))
				assert.Equal(t, tt.expected.errDebug, errutil.E(err).GetDebug())
			}
			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func requestWithBody(body io.ReadCloser) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Body = body
	return request
}

type errorReadCloser struct {
	err error
}

func (r errorReadCloser) Read(_ []byte) (int, error) {
	return 0, r.err
}

func (errorReadCloser) Close() error {
	return nil
}
