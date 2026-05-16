package dm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func TestService_createDMMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dmID := domain.DMID(10)
	messageID := domain.DMMessageID(20)
	senderID := domain.UserID(1)
	recipientID := domain.UserID(2)
	participants := []domain.UserID{senderID, recipientID}
	sender := &domain.User{
		Username: "alice",
		ID:       senderID,
	}

	type args struct {
		ctx     context.Context
		message *domain.DMMessage
	}
	type expected struct {
		err       expectedError
		messageID domain.DMMessageID
		hasTime   bool
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
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					CreateDMMessage(a.ctx, a.message).
					Return(messageID, nil)
				f.db.EXPECT().
					SetLastMessageInDM(a.ctx, a.message.DMID, messageID).
					Return(nil)
				f.db.EXPECT().
					SetDMLastReadMessage(a.ctx, a.message.DMID, a.message.SenderID, messageID).
					Return(nil)
				f.user.EXPECT().
					GetByID(a.message.SenderID).
					Return(sender, nil)
				f.db.EXPECT().
					GetDMParticipants(a.ctx, a.message.DMID).
					Return(participants, nil)
				f.redis.EXPECT().
					PublishToUser(a.ctx, senderID, gomock.Any()).
					Return(nil)
				f.redis.EXPECT().
					PublishToUser(a.ctx, recipientID, gomock.Any()).
					Return(nil)
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					Text:     "hello",
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				messageID: messageID,
				hasTime:   true,
			},
		},
		{
			name: "participant_not_found",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(nil, notFoundError())
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				err: wrappedError(errutil.NotFound, "user id='1' not a participant of dm id='10'"),
			},
		},
		{
			name: "create_message_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					CreateDMMessage(a.ctx, a.message).
					Return(domain.DMMessageID(0), assert.AnError)
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				err: dependencyError("s.db.CreateDMMessage"),
			},
		},
		{
			name: "set_last_message_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					CreateDMMessage(a.ctx, a.message).
					Return(messageID, nil)
				f.db.EXPECT().
					SetLastMessageInDM(a.ctx, a.message.DMID, messageID).
					Return(assert.AnError)
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				err: dependencyError("s.db.SetLastMessageInDM"),
			},
		},
		{
			name: "set_last_read_message_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					CreateDMMessage(a.ctx, a.message).
					Return(messageID, nil)
				f.db.EXPECT().
					SetLastMessageInDM(a.ctx, a.message.DMID, messageID).
					Return(nil)
				f.db.EXPECT().
					SetDMLastReadMessage(a.ctx, a.message.DMID, a.message.SenderID, messageID).
					Return(assert.AnError)
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				err: dependencyError("s.db.SetDMLastReadMessage"),
			},
		},
		{
			name: "get_participants_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					CreateDMMessage(a.ctx, a.message).
					Return(messageID, nil)
				f.db.EXPECT().
					SetLastMessageInDM(a.ctx, a.message.DMID, messageID).
					Return(nil)
				f.db.EXPECT().
					SetDMLastReadMessage(a.ctx, a.message.DMID, a.message.SenderID, messageID).
					Return(nil)
				f.user.EXPECT().
					GetByID(a.message.SenderID).
					Return(sender, nil)
				f.db.EXPECT().
					GetDMParticipants(a.ctx, a.message.DMID).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				err:       dependencyError("s.db.GetDMParticipants"),
				messageID: messageID,
				hasTime:   true,
			},
		},
		{
			name: "publish_error_is_ignored",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.message.DMID, a.message.SenderID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					CreateDMMessage(a.ctx, a.message).
					Return(messageID, nil)
				f.db.EXPECT().
					SetLastMessageInDM(a.ctx, a.message.DMID, messageID).
					Return(nil)
				f.db.EXPECT().
					SetDMLastReadMessage(a.ctx, a.message.DMID, a.message.SenderID, messageID).
					Return(nil)
				f.user.EXPECT().
					GetByID(a.message.SenderID).
					Return(sender, nil)
				f.db.EXPECT().
					GetDMParticipants(a.ctx, a.message.DMID).
					Return([]domain.UserID{senderID}, nil)
				f.redis.EXPECT().
					PublishToUser(a.ctx, senderID, gomock.Any()).
					Return(assert.AnError)
			},
			args: args{
				ctx: ctx,
				message: &domain.DMMessage{
					DMID:     dmID,
					SenderID: senderID,
				},
			},
			expected: expected{
				messageID: messageID,
				hasTime:   true,
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

			err := fields.service.createDMMessage(tt.args.ctx, tt.args.message)

			requireExpectedError(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.messageID, tt.args.message.ID)
			assert.Equal(t, tt.expected.hasTime, !tt.args.message.CreatedAt.IsZero())
		})
	}
}

func TestService_convertMessageToDesc(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2024, time.January, 2, 15, 4, 0, 0, time.UTC)
	senderID := domain.UserID(1)
	message := &domain.DMMessage{
		CreatedAt: createdAt,
		Text:      "hello",
		ID:        20,
		DMID:      10,
		SenderID:  senderID,
	}
	sender := &domain.User{
		Username: "alice",
		ID:       senderID,
	}

	type expected struct {
		payload *event_desc.MessagePayload
		channel event_desc.Channel
	}
	tests := []struct {
		prepare  func(*serviceFields)
		name     string
		expected expected
	}{
		{
			name: "success",
			prepare: func(f *serviceFields) {
				f.user.EXPECT().
					GetByID(senderID).
					Return(sender, nil)
			},
			expected: expected{
				channel: event_desc.Channel{
					Type: event_desc.ChannelDM,
					ID:   10,
				},
				payload: &event_desc.MessagePayload{
					ID:        20,
					Text:      "hello",
					CreatedAt: "Jan 2, 15:04",
					Sender: event_desc.UserInfo{
						ID:       1,
						Username: "alice",
						Avatar:   "/v1/image/alice.jpeg",
					},
				},
			},
		},
		{
			name: "unknown_user",
			prepare: func(f *serviceFields) {
				f.user.EXPECT().
					GetByID(senderID).
					Return(nil, assert.AnError)
			},
			expected: expected{
				channel: event_desc.Channel{
					Type: event_desc.ChannelDM,
					ID:   10,
				},
				payload: &event_desc.MessagePayload{
					ID:        20,
					Text:      "hello",
					CreatedAt: "Jan 2, 15:04",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			fields := newServiceFields(ctrl, 2)
			if tt.prepare != nil {
				tt.prepare(fields)
			}

			channel, payload := fields.service.convertMessageToDesc(message)

			require.Equal(t, tt.expected.channel, channel)
			assert.Equal(t, tt.expected.payload, payload)
		})
	}
}
