package bind

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chattery/internal/domain"
)

const invalidEnvValue = "invalid"

func TestEnvString(t *testing.T) {
	t.Parallel()

	configured := "configured"

	type args struct {
		envName      string
		defaultValue string
	}
	type expected struct {
		value string
	}
	tests := []struct {
		name     string
		envValue *string
		args     args
		expected expected
	}{
		{
			name: "default",
			args: args{
				envName:      "CHATTERY_TEST_ENV_STRING_DEFAULT",
				defaultValue: "default",
			},
			expected: expected{
				value: "default",
			},
		},
		{
			name:     "configured",
			envValue: &configured,
			args: args{
				envName:      "CHATTERY_TEST_ENV_STRING_CONFIGURED",
				defaultValue: "default",
			},
			expected: expected{
				value: configured,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prepareEnv(t, tt.args.envName, tt.envValue)

			value := EnvString(tt.args.envName, tt.args.defaultValue)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestEnvStrings(t *testing.T) {
	t.Parallel()

	configured := " first, ,second, third "

	type args struct {
		envName      string
		defaultValue []string
	}
	type expected struct {
		value []string
	}
	tests := []struct {
		name     string
		envValue *string
		args     args
		expected expected
	}{
		{
			name: "default",
			args: args{
				envName:      "CHATTERY_TEST_ENV_STRINGS_DEFAULT",
				defaultValue: []string{"default"},
			},
			expected: expected{
				value: []string{"default"},
			},
		},
		{
			name:     "configured",
			envValue: &configured,
			args: args{
				envName:      "CHATTERY_TEST_ENV_STRINGS_CONFIGURED",
				defaultValue: []string{"default"},
			},
			expected: expected{
				value: []string{"first", "second", "third"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prepareEnv(t, tt.args.envName, tt.envValue)

			value := EnvStrings(tt.args.envName, tt.args.defaultValue)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestEnvInt(t *testing.T) {
	t.Parallel()

	configured := "42"
	invalid := invalidEnvValue

	type args struct {
		envName      string
		defaultValue int
	}
	type expected struct {
		value int
	}
	tests := []struct {
		name     string
		envValue *string
		args     args
		expected expected
	}{
		{
			name: "default",
			args: args{
				envName:      "CHATTERY_TEST_ENV_INT_DEFAULT",
				defaultValue: 10,
			},
			expected: expected{
				value: 10,
			},
		},
		{
			name:     "configured",
			envValue: &configured,
			args: args{
				envName:      "CHATTERY_TEST_ENV_INT_CONFIGURED",
				defaultValue: 10,
			},
			expected: expected{
				value: 42,
			},
		},
		{
			name:     "invalid",
			envValue: &invalid,
			args: args{
				envName:      "CHATTERY_TEST_ENV_INT_INVALID",
				defaultValue: 10,
			},
			expected: expected{
				value: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prepareEnv(t, tt.args.envName, tt.envValue)

			value := EnvInt(tt.args.envName, tt.args.defaultValue)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestEnvInt64(t *testing.T) {
	t.Parallel()

	configured := "42"
	invalid := invalidEnvValue

	type args struct {
		envName      string
		defaultValue int64
	}
	type expected struct {
		value int64
	}
	tests := []struct {
		name     string
		envValue *string
		args     args
		expected expected
	}{
		{
			name: "default",
			args: args{
				envName:      "CHATTERY_TEST_ENV_INT64_DEFAULT",
				defaultValue: 10,
			},
			expected: expected{
				value: 10,
			},
		},
		{
			name:     "configured",
			envValue: &configured,
			args: args{
				envName:      "CHATTERY_TEST_ENV_INT64_CONFIGURED",
				defaultValue: 10,
			},
			expected: expected{
				value: 42,
			},
		},
		{
			name:     "invalid",
			envValue: &invalid,
			args: args{
				envName:      "CHATTERY_TEST_ENV_INT64_INVALID",
				defaultValue: 10,
			},
			expected: expected{
				value: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prepareEnv(t, tt.args.envName, tt.envValue)

			value := EnvInt64(tt.args.envName, tt.args.defaultValue)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestEnvDuration(t *testing.T) {
	t.Parallel()

	configured := "2m"
	invalid := invalidEnvValue

	type args struct {
		envName      string
		defaultValue time.Duration
	}
	type expected struct {
		value time.Duration
	}
	tests := []struct {
		name     string
		envValue *string
		args     args
		expected expected
	}{
		{
			name: "default",
			args: args{
				envName:      "CHATTERY_TEST_ENV_DURATION_DEFAULT",
				defaultValue: 5 * time.Second,
			},
			expected: expected{
				value: 5 * time.Second,
			},
		},
		{
			name:     "configured",
			envValue: &configured,
			args: args{
				envName:      "CHATTERY_TEST_ENV_DURATION_CONFIGURED",
				defaultValue: 5 * time.Second,
			},
			expected: expected{
				value: 2 * time.Minute,
			},
		},
		{
			name:     "invalid",
			envValue: &invalid,
			args: args{
				envName:      "CHATTERY_TEST_ENV_DURATION_INVALID",
				defaultValue: 5 * time.Second,
			},
			expected: expected{
				value: 5 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prepareEnv(t, tt.args.envName, tt.envValue)

			value := EnvDuration(tt.args.envName, tt.args.defaultValue)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestEnvBool(t *testing.T) {
	t.Parallel()

	configured := "true"
	invalid := invalidEnvValue

	type args struct {
		envName      string
		defaultValue bool
	}
	type expected struct {
		value bool
	}
	tests := []struct {
		name     string
		envValue *string
		args     args
		expected expected
	}{
		{
			name: "default",
			args: args{
				envName:      "CHATTERY_TEST_ENV_BOOL_DEFAULT",
				defaultValue: true,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name:     "configured",
			envValue: &configured,
			args: args{
				envName:      "CHATTERY_TEST_ENV_BOOL_CONFIGURED",
				defaultValue: false,
			},
			expected: expected{
				value: true,
			},
		},
		{
			name:     "invalid",
			envValue: &invalid,
			args: args{
				envName:      "CHATTERY_TEST_ENV_BOOL_INVALID",
				defaultValue: false,
			},
			expected: expected{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			prepareEnv(t, tt.args.envName, tt.envValue)

			value := EnvBool(tt.args.envName, tt.args.defaultValue)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestPathParamI64(t *testing.T) {
	t.Parallel()

	type args struct {
		r         *http.Request
		paramName string
	}
	type expected struct {
		value domain.UserID
		err   bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				r:         requestWithPathParam("user_id", "42"),
				paramName: "user_id",
			},
			expected: expected{
				value: 42,
			},
		},
		{
			name: "invalid",
			args: args{
				r:         requestWithPathParam("user_id", "invalid"),
				paramName: "user_id",
			},
			expected: expected{
				value: domain.UserIsUnknown,
				err:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := PathParamI64[domain.UserID](tt.args.r, tt.args.paramName)

			if tt.expected.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func prepareEnv(t *testing.T, envName string, envValue *string) {
	t.Helper()

	if envValue == nil {
		require.NoError(t, os.Unsetenv(envName))
		return
	}

	require.NoError(t, os.Setenv(envName, *envValue))
	t.Cleanup(func() {
		require.NoError(t, os.Unsetenv(envName))
	})
}

func requestWithPathParam(name, value string) *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(name, value)
	return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
}
