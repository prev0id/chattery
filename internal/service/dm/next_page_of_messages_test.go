package dm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func TestService_nextPagesOfDMMessages(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	limit := 2
	dmID := domain.DMID(10)
	userID := domain.UserID(1)
	firstCreatedAt := time.Date(2024, time.January, 3, 10, 0, 0, 0, time.UTC)
	secondCreatedAt := time.Date(2024, time.January, 2, 10, 0, 0, 0, time.UTC)
	firstMessage := &domain.DMMessage{
		CreatedAt: firstCreatedAt,
		ID:        30,
	}
	secondMessage := &domain.DMMessage{
		CreatedAt: secondCreatedAt,
		ID:        20,
	}

	type args struct {
		ctx    context.Context
		cursor *domain.DMCursor
		userID domain.UserID
	}
	type expected struct {
		nextCursor *domain.DMCursor
		messages   []*domain.DMMessage
		err        expectedError
		limit      int
	}
	tests := []struct {
		prepare  func(*serviceFields, *args)
		name     string
		args     args
		expected expected
	}{
		{
			name: "success_with_next_cursor",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.cursor.ChatID, a.userID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					NextPagesOfDMMessages(a.ctx, a.cursor).
					Return([]*domain.DMMessage{firstMessage, secondMessage}, nil)
			},
			args: args{
				ctx: ctx,
				cursor: &domain.DMCursor{
					ChatID: dmID,
				},
				userID: userID,
			},
			expected: expected{
				messages: []*domain.DMMessage{firstMessage, secondMessage},
				nextCursor: &domain.DMCursor{
					ChatID:    dmID,
					MessageID: secondMessage.ID,
					Timestamp: secondCreatedAt,
				},
				limit: limit,
			},
		},
		{
			name: "success_without_next_cursor",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.cursor.ChatID, a.userID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					NextPagesOfDMMessages(a.ctx, a.cursor).
					Return([]*domain.DMMessage{firstMessage}, nil)
			},
			args: args{
				ctx: ctx,
				cursor: &domain.DMCursor{
					ChatID: dmID,
				},
				userID: userID,
			},
			expected: expected{
				messages: []*domain.DMMessage{firstMessage},
				limit:    limit,
			},
		},
		{
			name: "participant_not_found",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.cursor.ChatID, a.userID).
					Return(nil, notFoundError())
			},
			args: args{
				ctx: ctx,
				cursor: &domain.DMCursor{
					ChatID: dmID,
				},
				userID: userID,
			},
			expected: expected{
				err: wrappedError(errutil.NotFound, "user id='1' not a participant of dm id='10'"),
			},
		},
		{
			name: "list_messages_error",
			prepare: func(f *serviceFields, a *args) {
				f.db.EXPECT().
					GetDMParticipant(a.ctx, a.cursor.ChatID, a.userID).
					Return(&domain.DMParticipant{}, nil)
				f.db.EXPECT().
					NextPagesOfDMMessages(a.ctx, a.cursor).
					Return(nil, assert.AnError)
			},
			args: args{
				ctx: ctx,
				cursor: &domain.DMCursor{
					ChatID: dmID,
				},
				userID: userID,
			},
			expected: expected{
				err:   dependencyError("s.db.NextPagesOfDMMessages"),
				limit: limit,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			fields := newServiceFields(ctrl, limit)
			if tt.prepare != nil {
				tt.prepare(fields, &tt.args)
			}

			messages, nextCursor, err := fields.service.nextPagesOfDMMessages(
				tt.args.ctx,
				tt.args.userID,
				tt.args.cursor,
			)

			requireExpectedError(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.messages, messages)
			assert.Equal(t, tt.expected.nextCursor, nextCursor)
			assert.Equal(t, tt.expected.limit, tt.args.cursor.Limit)
		})
	}
}
