package dm

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"chattery/internal/domain"
	"chattery/internal/service/dm/mock_dm"
	"chattery/internal/utils/errutil"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

type serviceFields struct {
	db          *mock_dm.Mockdb
	redis       *mock_dm.Mockredis
	user        *mock_dm.MockuserCache
	cache       *mock_dm.MockdmCache
	transaction *mock_dm.MocktxManager
	service     *Service
}

type expectedError struct {
	err     error
	message string
	debug   []string
	kind    errutil.Kind
	hasErr  bool
}

func newServiceFields(ctrl *gomock.Controller, limit int) *serviceFields {
	fields := &serviceFields{
		db:          mock_dm.NewMockdb(ctrl),
		redis:       mock_dm.NewMockredis(ctrl),
		user:        mock_dm.NewMockuserCache(ctrl),
		cache:       mock_dm.NewMockdmCache(ctrl),
		transaction: mock_dm.NewMocktxManager(ctrl),
	}
	fields.service = &Service{
		db:          fields.db,
		transaction: fields.transaction,
		redis:       fields.redis,
		user:        fields.user,
		cache:       fields.cache,
		limit:       limit,
	}
	return fields
}

func notFoundError() error {
	return errutil.E(assert.AnError).Kind(errutil.NotFound)
}

func businessError(kind errutil.Kind, message string) expectedError {
	return expectedError{
		kind:    kind,
		message: message,
		hasErr:  true,
	}
}

func dependencyError(debug string) expectedError {
	return expectedError{
		err:    assert.AnError,
		debug:  []string{debug},
		kind:   errutil.Internal,
		hasErr: true,
	}
}

func wrappedError(kind errutil.Kind, message string) expectedError {
	return expectedError{
		err:     assert.AnError,
		message: message,
		kind:    kind,
		hasErr:  true,
	}
}

func requireExpectedError(t *testing.T, actual error, expected expectedError) {
	t.Helper()

	if !expected.hasErr {
		require.NoError(t, actual)
		return
	}

	require.Error(t, actual)
	assert.Equal(t, expected.err, errutil.E(actual).GetError())
	assert.Equal(t, expected.debug, errutil.E(actual).GetDebug())
	assert.Equal(t, expected.message, errutil.E(actual).GetMessage())
	assert.Equal(t, expected.kind, errutil.E(actual).GetKind())
}

func TestService_getNextCursor(t *testing.T) {
	t.Parallel()

	dmID := domain.DMID(10)
	createdAt := time.Date(2024, time.January, 2, 10, 0, 0, 0, time.UTC)
	message := &domain.DMMessage{
		CreatedAt: createdAt,
		ID:        20,
	}

	type fields struct {
		service *Service
	}
	type args struct {
		messages []*domain.DMMessage
		dmID     domain.DMID
	}
	type expected struct {
		cursor *domain.DMCursor
	}
	tests := []struct {
		fields   fields
		expected expected
		name     string
		args     args
	}{
		{
			name: "full_page",
			fields: fields{
				service: &Service{
					limit: 1,
				},
			},
			args: args{
				dmID:     dmID,
				messages: []*domain.DMMessage{message},
			},
			expected: expected{
				cursor: &domain.DMCursor{
					ChatID:    dmID,
					MessageID: message.ID,
					Timestamp: createdAt,
				},
			},
		},
		{
			name: "not_full_page",
			fields: fields{
				service: &Service{
					limit: 2,
				},
			},
			args: args{
				dmID:     dmID,
				messages: []*domain.DMMessage{message},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.fields.service.getNextCursor(tt.args.dmID, tt.args.messages)

			assert.Equal(t, tt.expected.cursor, value)
		})
	}
}
