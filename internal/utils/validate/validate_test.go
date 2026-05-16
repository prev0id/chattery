package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/validate"
)

func TestUsername(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	tests := []struct {
		name     string
		username string
		expected expected
	}{
		{
			name:     "valid",
			username: "user_1-name",
		},
		{
			name:     "empty",
			username: "",
			expected: expected{
				errMessage: "username must be provided",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:     "too_long",
			username: "abcdefghijklmnopqrstuvwxyz",
			expected: expected{
				errMessage: "username must be at most 25 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:     "invalid_characters",
			username: "user name",
			expected: expected{
				errMessage: "username can only contain letters (a-z, A-Z), digits, underscores and dashes",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.Username(tt.username)

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	tests := []struct {
		name     string
		login    string
		expected expected
	}{
		{
			name:  "valid",
			login: "user@example.com",
		},
		{
			name:  "invalid",
			login: "plain-address",
			expected: expected{
				errMessage: "login must be a valid email address",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.Login(tt.login)

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func TestPassword(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	tests := []struct {
		name     string
		password string
		expected expected
	}{
		{
			name:     "valid",
			password: "Password1",
		},
		{
			name:     "too_short",
			password: "Pass1",
			expected: expected{
				errMessage: "password must be at least 8 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:     "too_long",
			password: "Password1Password1Password1Password1",
			expected: expected{
				errMessage: "password must be at most 32 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:     "without_lowercase_letter",
			password: "PASSWORD1",
			expected: expected{
				errMessage: "password must contain at least one lowercase letter",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:     "without_uppercase_letter",
			password: "password1",
			expected: expected{
				errMessage: "password must contain at least one uppercase letter",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:     "without_digit",
			password: "Password",
			expected: expected{
				errMessage: "password must contain at least one digit",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.Password(tt.password)

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func TestServerName(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	tests := []struct {
		name       string
		serverName string
		expected   expected
	}{
		{
			name:       "valid",
			serverName: "Server_1-Name",
		},
		{
			name:       "too_short",
			serverName: "srv",
			expected: expected{
				errMessage: "name must be at least 5 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:       "too_long",
			serverName: "abcdefghijklmnopqrstuvwxyz",
			expected: expected{
				errMessage: "name must be at most 25 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:       "invalid_characters",
			serverName: "server!",
			expected: expected{
				errMessage: "name can only contain letters (a-z, A-Z), digits, spaces, underscores and dashes -",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.ServerName(tt.serverName)

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func TestTopicName(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	tests := []struct {
		name      string
		topicName string
		expected  expected
	}{
		{
			name:      "valid",
			topicName: "Topic 1",
		},
		{
			name:      "too_short",
			topicName: "t",
			expected: expected{
				errMessage: "name must be at least 2 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:      "too_long",
			topicName: "abcdefghijklmnopqrstu",
			expected: expected{
				errMessage: "name must be at most 20 characters long",
				errKind:    errutil.InvalidRequest,
			},
		},
		{
			name:      "invalid_characters",
			topicName: "topic!",
			expected: expected{
				errMessage: "name can only contain letters (a-z, A-Z), digits, spaces, underscores and dashes -",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.TopicName(tt.topicName)

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func TestTopicType(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	tests := []struct {
		name      string
		topicType string
		expected  expected
	}{
		{
			name:      "text",
			topicType: domain.TopicTypeText.String(),
		},
		{
			name:      "voice",
			topicType: domain.TopicTypeVoice.String(),
		},
		{
			name:      "invalid",
			topicType: "video",
			expected: expected{
				errMessage: "type must be one of values [text voice]",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.TopicType(tt.topicType)

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func TestNotEmpty(t *testing.T) {
	t.Parallel()

	type expected struct {
		errMessage string
		errKind    errutil.Kind
	}

	stringTests := []struct {
		name     string
		value    string
		expected expected
	}{
		{
			name:  "string_value_is_valid",
			value: "value",
		},
		{
			name:  "zero_string_is_invalid",
			value: "",
			expected: expected{
				errMessage: "field must be provided",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.NotEmpty(tt.value, "field")

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}

	intTests := []struct {
		name     string
		expected expected
		value    int
	}{
		{
			name:  "int_value_is_valid",
			value: 1,
		},
		{
			name:  "zero_int_is_invalid",
			value: 0,
			expected: expected{
				errMessage: "field must be provided",
				errKind:    errutil.InvalidRequest,
			},
		},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate.NotEmpty(tt.value, "field")

			requireValidationResult(t, err, tt.expected.errKind, tt.expected.errMessage)
		})
	}
}

func requireValidationResult(t *testing.T, err error, wantKind errutil.Kind, wantMessage string) {
	t.Helper()

	if wantMessage == "" {
		require.NoError(t, err)
		return
	}

	require.Error(t, err)
	assert.True(t, errutil.Is(wantKind, err))
	assert.Equal(t, wantMessage, errutil.E(err).GetMessage())
}
