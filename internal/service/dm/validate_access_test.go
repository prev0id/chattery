package dm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
)

func TestService_ValidateAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := domain.UserID(1)
	dmID := domain.DMID(10)

	type args struct {
		ctx    context.Context
		userID domain.UserID
		dmID   domain.DMID
	}
	type expected struct {
		err error
	}
	tests := []struct {
		prepare  func(*serviceFields, *args)
		name     string
		expected expected
		args     args
	}{
		{
			name: "success",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.dmID, a.userID).
					Return(&domain.DMParticipant{}, nil)
			},
			args: args{
				ctx:    ctx,
				userID: userID,
				dmID:   dmID,
			},
		},
		{
			name: "error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.dmID, a.userID).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx:    ctx,
				userID: userID,
				dmID:   dmID,
			},
			expected: expected{
				err: assert.AnError,
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

			err := fields.service.ValidateAccess(tt.args.ctx, tt.args.userID, tt.args.dmID)

			assert.Equal(t, tt.expected.err, err)
		})
	}
}
