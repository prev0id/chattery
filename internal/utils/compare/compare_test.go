package compare

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"chattery/internal/domain"
)

func TestServers(t *testing.T) {
	t.Parallel()

	joinedAt := time.Date(2026, time.May, 1, 12, 0, 0, 0, time.UTC)

	type args struct {
		lhs *domain.Server
		rhs *domain.Server
	}
	type expected struct {
		value int
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "less_by_joined_at",
			args: args{
				lhs: &domain.Server{ID: 2, JoinedAt: joinedAt},
				rhs: &domain.Server{ID: 1, JoinedAt: joinedAt.Add(time.Hour)},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "less_by_id",
			args: args{
				lhs: &domain.Server{ID: 1, JoinedAt: joinedAt},
				rhs: &domain.Server{ID: 2, JoinedAt: joinedAt},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "equal",
			args: args{
				lhs: &domain.Server{ID: 1, JoinedAt: joinedAt},
				rhs: &domain.Server{ID: 1, JoinedAt: joinedAt},
			},
			expected: expected{
				value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := Servers(tt.args.lhs, tt.args.rhs)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestServersByName(t *testing.T) {
	t.Parallel()

	type args struct {
		lhs *domain.Server
		rhs *domain.Server
	}
	type expected struct {
		value int
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "less_by_name",
			args: args{
				lhs: &domain.Server{ID: 2, Name: "alpha"},
				rhs: &domain.Server{ID: 1, Name: "beta"},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "less_by_id",
			args: args{
				lhs: &domain.Server{ID: 1, Name: "alpha"},
				rhs: &domain.Server{ID: 2, Name: "alpha"},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "equal",
			args: args{
				lhs: &domain.Server{ID: 1, Name: "alpha"},
				rhs: &domain.Server{ID: 1, Name: "alpha"},
			},
			expected: expected{
				value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := ServersByName(tt.args.lhs, tt.args.rhs)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestUsersByUsername(t *testing.T) {
	t.Parallel()

	type args struct {
		lhs *domain.User
		rhs *domain.User
	}
	type expected struct {
		value int
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "less_by_username",
			args: args{
				lhs: &domain.User{ID: 2, Username: "alice"},
				rhs: &domain.User{ID: 1, Username: "bob"},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "less_by_id",
			args: args{
				lhs: &domain.User{ID: 1, Username: "alice"},
				rhs: &domain.User{ID: 2, Username: "alice"},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "equal",
			args: args{
				lhs: &domain.User{ID: 1, Username: "alice"},
				rhs: &domain.User{ID: 1, Username: "alice"},
			},
			expected: expected{
				value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := UsersByUsername(tt.args.lhs, tt.args.rhs)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestTopics(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, time.May, 1, 12, 0, 0, 0, time.UTC)

	type args struct {
		lhs *domain.Topic
		rhs *domain.Topic
	}
	type expected struct {
		value int
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "less_by_created_at",
			args: args{
				lhs: &domain.Topic{ID: 2, CreatedAt: createdAt},
				rhs: &domain.Topic{ID: 1, CreatedAt: createdAt.Add(time.Hour)},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "less_by_id",
			args: args{
				lhs: &domain.Topic{ID: 1, CreatedAt: createdAt},
				rhs: &domain.Topic{ID: 2, CreatedAt: createdAt},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "equal",
			args: args{
				lhs: &domain.Topic{ID: 1, CreatedAt: createdAt},
				rhs: &domain.Topic{ID: 1, CreatedAt: createdAt},
			},
			expected: expected{
				value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := Topics(tt.args.lhs, tt.args.rhs)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestDMs(t *testing.T) {
	t.Parallel()

	lastActivityAt := time.Date(2026, time.May, 1, 12, 0, 0, 0, time.UTC)

	type args struct {
		lhs *domain.DM
		rhs *domain.DM
	}
	type expected struct {
		value int
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "less_by_last_activity_at",
			args: args{
				lhs: &domain.DM{ID: 2, LastActivityAt: lastActivityAt},
				rhs: &domain.DM{ID: 1, LastActivityAt: lastActivityAt.Add(time.Hour)},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "less_by_id",
			args: args{
				lhs: &domain.DM{ID: 1, LastActivityAt: lastActivityAt},
				rhs: &domain.DM{ID: 2, LastActivityAt: lastActivityAt},
			},
			expected: expected{
				value: -1,
			},
		},
		{
			name: "equal",
			args: args{
				lhs: &domain.DM{ID: 1, LastActivityAt: lastActivityAt},
				rhs: &domain.DM{ID: 1, LastActivityAt: lastActivityAt},
			},
			expected: expected{
				value: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := DMs(tt.args.lhs, tt.args.rhs)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}
