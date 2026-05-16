package dm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func TestService_createDM(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dmID := domain.DMID(10)
	userID1 := domain.UserID(1)
	userID2 := domain.UserID(2)

	type args struct {
		ctx          context.Context
		participant1 domain.UserID
		participant2 domain.UserID
	}
	type expected struct {
		err  expectedError
		dmID domain.DMID
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
					GetDMBetweenUsers(a.ctx, a.participant1, a.participant2).
					Return(nil, notFoundError())
				f.db.EXPECT().
					CreateDM(a.ctx).
					Return(dmID, nil)
				f.db.EXPECT().
					CreateDMParticipant(a.ctx, dmID, a.participant1).
					Return(nil)
				f.db.EXPECT().
					CreateDMParticipant(a.ctx, dmID, a.participant2).
					Return(nil)
			},
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID2,
			},
			expected: expected{
				dmID: dmID,
			},
		},
		{
			name: "same_user",
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID1,
			},
			expected: expected{
				err: businessError(errutil.InvalidRequest, "cannot create dm with yourself"),
			},
		},
		{
			name: "dm_already_exists",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMBetweenUsers(a.ctx, a.participant1, a.participant2).
					Return(&domain.DM{ID: dmID}, nil)
			},
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID2,
			},
			expected: expected{
				err: businessError(errutil.Exist, "dm already exists between users id='1' and id='2'"),
			},
		},
		{
			name: "get_dm_between_users_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMBetweenUsers(a.ctx, a.participant1, a.participant2).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID2,
			},
			expected: expected{
				err: dependencyError("s.db.GetDMBetweenUsers"),
			},
		},
		{
			name: "create_dm_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMBetweenUsers(a.ctx, a.participant1, a.participant2).
					Return(nil, notFoundError())
				f.db.EXPECT().
					CreateDM(a.ctx).
					Return(domain.DMID(0), assert.AnError)
			},
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID2,
			},
			expected: expected{
				err: dependencyError("s.db.CreateDM"),
			},
		},
		{
			name: "create_first_participant_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMBetweenUsers(a.ctx, a.participant1, a.participant2).
					Return(nil, notFoundError())
				f.db.EXPECT().
					CreateDM(a.ctx).
					Return(dmID, nil)
				f.db.EXPECT().
					CreateDMParticipant(a.ctx, dmID, a.participant1).
					Return(assert.AnError)
			},
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID2,
			},
			expected: expected{
				err: dependencyError("s.db.CreateDMParticipant"),
			},
		},
		{
			name: "create_second_participant_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMBetweenUsers(a.ctx, a.participant1, a.participant2).
					Return(nil, notFoundError())
				f.db.EXPECT().
					CreateDM(a.ctx).
					Return(dmID, nil)
				f.db.EXPECT().
					CreateDMParticipant(a.ctx, dmID, a.participant1).
					Return(nil)
				f.db.EXPECT().
					CreateDMParticipant(a.ctx, dmID, a.participant2).
					Return(assert.AnError)
			},
			args: args{
				ctx:          ctx,
				participant1: userID1,
				participant2: userID2,
			},
			expected: expected{
				err: dependencyError("s.db.CreateDMParticipant"),
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

			value, err := fields.service.createDM(tt.args.ctx, tt.args.participant1, tt.args.participant2)

			requireExpectedError(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.dmID, value)
		})
	}
}
