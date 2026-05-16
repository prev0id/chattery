package bind

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chattery/internal/utils/errutil"
)

type jsonTestValue struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func TestJSONString(t *testing.T) {
	t.Parallel()

	validValue := &jsonTestValue{Name: "alice", ID: 7}

	type args struct {
		raw string
	}
	type expected struct {
		value      *jsonTestValue
		errMessage string
		errDebug   []string
		errKind    errutil.Kind
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				raw: `{"name":"alice","id":7}`,
			},
			expected: expected{
				value: validValue,
			},
		},
		{
			name: "invalid_json",
			args: args{
				raw: `{"name":`,
			},
			expected: expected{
				errKind:    errutil.InvalidRequest,
				errMessage: "invalid json provided",
				errDebug:   []string{"json.Unmarshal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := JSONString[jsonTestValue](tt.args.raw)

			requireErrutilError(t, err, tt.expected.errKind, tt.expected.errMessage, tt.expected.errDebug)
			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestJSONBytes(t *testing.T) {
	t.Parallel()

	validValue := &jsonTestValue{Name: "alice", ID: 7}

	type args struct {
		raw []byte
	}
	type expected struct {
		value      *jsonTestValue
		errMessage string
		errDebug   []string
		errKind    errutil.Kind
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid",
			args: args{
				raw: []byte(`{"name":"alice","id":7}`),
			},
			expected: expected{
				value: validValue,
			},
		},
		{
			name: "invalid_json",
			args: args{
				raw: []byte(`{"name":`),
			},
			expected: expected{
				errKind:    errutil.InvalidRequest,
				errMessage: "invalid json provided",
				errDebug:   []string{"json.Unmarshal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := JSONBytes[jsonTestValue](tt.args.raw)

			requireErrutilError(t, err, tt.expected.errKind, tt.expected.errMessage, tt.expected.errDebug)
			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func requireErrutilError(
	t *testing.T,
	err error,
	expectedKind errutil.Kind,
	expectedMessage string,
	expectedDebug []string,
) {
	t.Helper()

	if expectedMessage == "" {
		require.NoError(t, err)
		return
	}

	require.Error(t, err)
	assert.True(t, errutil.Is(expectedKind, err))
	assert.Equal(t, expectedMessage, errutil.E(err).GetMessage())
	assert.Equal(t, expectedDebug, errutil.E(err).GetDebug())
}
