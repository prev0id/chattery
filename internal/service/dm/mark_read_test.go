package dm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func TestService_markDMMessageRead(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dmID := domain.DMID(10)
	userID := domain.UserID(1)
	messageID := domain.DMMessageID(20)

	type args struct {
		ctx       context.Context
		dmID      domain.DMID
		userID    domain.UserID
		messageID domain.DMMessageID
	}
	type expected struct {
		err expectedError
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
				f.db.EXPECT().
					GetDMMessage(a.ctx, a.dmID, a.messageID).
					Return(&domain.DMMessage{}, nil)
				f.db.EXPECT().
					SetDMLastReadMessage(a.ctx, a.dmID, a.userID, a.messageID).
					Return(nil)
			},
			args: args{
				ctx:       ctx,
				dmID:      dmID,
				userID:    userID,
				messageID: messageID,
			},
		},
		{
			name: "participant_not_found",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.dmID, a.userID).
					Return(nil, notFoundError())
			},
			args: args{
				ctx:       ctx,
				dmID:      dmID,
				userID:    userID,
				messageID: messageID,
			},
			expected: expected{
				err: wrappedError(errutil.NotFound, "user id='1' not a participant of dm id='10'"),
			},
		},
		{
			name: "message_not_found",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.dmID, a.userID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					GetDMMessage(a.ctx, a.dmID, a.messageID).
					Return(nil, notFoundError())
			},
			args: args{
				ctx:       ctx,
				dmID:      dmID,
				userID:    userID,
				messageID: messageID,
			},
			expected: expected{
				err: wrappedError(errutil.InvalidRequest, "message id='20' not found in dm id='10'"),
			},
		},
		{
			name: "message_lookup_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.dmID, a.userID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					GetDMMessage(a.ctx, a.dmID, a.messageID).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx:       ctx,
				dmID:      dmID,
				userID:    userID,
				messageID: messageID,
			},
			expected: expected{
				err: dependencyError("s.db.GetDMMessage"),
			},
		},
		{
			name: "set_last_read_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.dmID, a.userID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					GetDMMessage(a.ctx, a.dmID, a.messageID).
					Return(&domain.DMMessage{}, nil)
				f.db.EXPECT().
					SetDMLastReadMessage(a.ctx, a.dmID, a.userID, a.messageID).
					Return(assert.AnError)
			},
			args: args{
				ctx:       ctx,
				dmID:      dmID,
				userID:    userID,
				messageID: messageID,
			},
			expected: expected{
				err: dependencyError("s.db.SetDMLastReadMessage"),
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

			err := fields.service.markDMMessageRead(tt.args.ctx, tt.args.userID, tt.args.dmID, tt.args.messageID)

			requireExpectedError(t, err, tt.expected.err)
		})
	}
}
