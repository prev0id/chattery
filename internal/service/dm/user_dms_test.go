package dm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
)

func TestService_UserDMs(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := domain.UserID(1)
	earlier := time.Date(2024, time.January, 1, 10, 0, 0, 0, time.UTC)
	later := time.Date(2024, time.January, 2, 10, 0, 0, 0, time.UTC)
	dm1 := &domain.DM{
		LastActivityAt: later,
		ID:             1,
	}
	dm2 := &domain.DM{
		LastActivityAt: earlier,
		ID:             2,
	}
	dm3 := &domain.DM{
		LastActivityAt: earlier,
		ID:             1,
	}

	type args struct {
		ctx    context.Context
		userID domain.UserID
	}
	type expected struct {
		dms []*domain.DM
		err expectedError
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
					UserDMs(a.ctx, a.userID).
					Return([]*domain.DM{dm1, dm2, dm3}, nil)
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			expected: expected{
				dms: []*domain.DM{dm3, dm2, dm1},
			},
		},
		{
			name: "db_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					UserDMs(a.ctx, a.userID).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			expected: expected{
				err: dependencyError("s.db.UserDMs"),
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

			dms, err := fields.service.UserDMs(tt.args.ctx, tt.args.userID)

			requireExpectedError(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.dms, dms)
		})
	}
}
