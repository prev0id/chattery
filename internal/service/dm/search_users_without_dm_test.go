package dm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
)

func TestService_SearchUsersWithoutDM(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := domain.UserID(1)
	alice := &domain.User{
		Username: "Alice",
		ID:       4,
	}
	bob := &domain.User{
		Username: "Bob",
		ID:       2,
	}
	mallory := &domain.User{
		Username: "Mallory",
		ID:       3,
	}
	self := &domain.User{
		Username: "Alina",
		ID:       userID,
	}
	users := []*domain.User{mallory, bob, alice, self}

	type args struct {
		ctx    context.Context
		query  string
		userID domain.UserID
	}
	type expected struct {
		users []*domain.User
		err   expectedError
	}
	tests := []struct {
		prepare  func(*serviceFields, *args)
		name     string
		args     args
		expected expected
	}{
		{
			name: "success",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetUserDMParticipantIDs(a.ctx, a.userID).
					Return([]domain.UserID{bob.ID}, nil)
				f.user.EXPECT().
					List().
					Return(users)
			},
			args: args{
				ctx:    ctx,
				query:  " AL ",
				userID: userID,
			},
			expected: expected{
				users: []*domain.User{alice, mallory},
			},
		},
		{
			name: "db_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetUserDMParticipantIDs(a.ctx, a.userID).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx:    ctx,
				query:  "al",
				userID: userID,
			},
			expected: expected{
				err: dependencyError("s.db.GetUserDMParticipantIDs"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			fields := newServiceFields(ctrl, 2)
			if tt.prepare != nil {
				tt.prepare(fields, &tt.args)
			}

			users, err := fields.service.SearchUsersWithoutDM(tt.args.ctx, tt.args.userID, tt.args.query)

			requireExpectedError(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.users, users)
		})
	}
}
