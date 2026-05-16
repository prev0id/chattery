package transaction

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"chattery/internal/client/postgres"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/transaction/mock_transaction"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestNewManager(t *testing.T) {
	t.Parallel()

	type args struct {
		pool *pgxpool.Pool
	}
	type expected struct {
		pool      *pgxpool.Pool
		hasClient bool
	}
	tests := []struct {
		args     args
		name     string
		expected expected
	}{
		{
			name: "nil_pool",
			expected: expected{
				hasClient: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := NewManager(tt.args.pool)

			assert.Equal(t, tt.expected.pool, value.pool)
			assert.Equal(t, tt.expected.hasClient, value.client != nil)
		})
	}
}

func TestManager_Query(t *testing.T) {
	t.Parallel()

	client := postgres.New(nil)

	type fields struct {
		manager *Manager
	}
	type args struct {
		ctx context.Context
	}
	type expected struct {
		value postgres.Querier
	}
	tests := []struct {
		fields   fields
		args     args
		expected expected
		name     string
	}{
		{
			name: "without_tx",
			fields: fields{
				manager: &Manager{
					client: client,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				value: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.fields.manager.Query(tt.args.ctx)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func TestManager_QueryWithTx(t *testing.T) {
	t.Parallel()

	client := postgres.New(nil)
	createUser := &postgres.CreateUserParams{
		Login:    "login",
		Password: []byte("password"),
		Username: "username",
	}
	createdUserID := int64(1)

	type fields struct {
		manager *Manager
		row     *mock_transaction.MockRow
		tx      *mock_transaction.MockTx
	}
	type args struct {
		createUser *postgres.CreateUserParams
		ctx        context.Context
	}
	type expected struct {
		err          error
		createUserID int64
	}
	tests := []struct {
		prepare  func(*fields, *args)
		fields   fields
		args     args
		expected expected
		name     string
	}{
		{
			name: "with_tx",
			prepare: func(f *fields, a *args) {
				a.ctx = txToContext(context.Background(), f.tx)
				a.createUser = createUser
				f.tx.EXPECT().
					QueryRow(a.ctx, gomock.Any(), createUser.Login, createUser.Password, createUser.Username).
					Return(f.row)
				f.row.EXPECT().
					Scan(gomock.Any()).
					DoAndReturn(scanInt64(createdUserID))
			},
			fields: fields{
				manager: &Manager{
					client: client,
				},
			},
			expected: expected{
				createUserID: createdUserID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields.row = mock_transaction.NewMockRow(ctrl)
			tt.fields.tx = mock_transaction.NewMockTx(ctrl)
			if tt.prepare != nil {
				tt.prepare(&tt.fields, &tt.args)
			}

			value := tt.fields.manager.Query(tt.args.ctx)
			id, err := value.CreateUser(tt.args.ctx, tt.args.createUser)

			requireError(t, err, tt.expected.err, nil)
			assert.Equal(t, tt.expected.createUserID, id)
		})
	}
}

func TestManager_InTransaction(t *testing.T) {
	t.Parallel()

	type fields struct {
		manager *Manager
		tx      *mock_transaction.MockTx
	}
	type args struct {
		ctx context.Context
		fn  func(*fields, context.Context) error
	}
	type expected struct {
		err error
	}
	tests := []struct {
		prepare  func(*fields, *args)
		fields   fields
		args     args
		expected expected
		name     string
	}{
		{
			name: "uses_existing_tx",
			prepare: func(f *fields, a *args) {
				a.ctx = txToContext(context.Background(), f.tx)
				a.fn = func(f *fields, ctx context.Context) error {
					if txFromContext(ctx) != f.tx {
						return assert.AnError
					}
					return nil
				}
			},
			fields: fields{
				manager: &Manager{},
			},
		},
		{
			name: "returns_fn_error_with_existing_tx",
			prepare: func(f *fields, a *args) {
				a.ctx = txToContext(context.Background(), f.tx)
				a.fn = func(f *fields, ctx context.Context) error {
					if txFromContext(ctx) != f.tx {
						return assert.AnError
					}
					return assert.AnError
				}
			},
			fields: fields{
				manager: &Manager{},
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
			tt.fields.tx = mock_transaction.NewMockTx(ctrl)
			tt.prepare(&tt.fields, &tt.args)

			err := tt.fields.manager.InTransaction(tt.args.ctx, func(ctx context.Context) error {
				return tt.args.fn(&tt.fields, ctx)
			})

			requireError(t, err, tt.expected.err, nil)
		})
	}
}

func Test_rollback(t *testing.T) {
	t.Parallel()

	type fields struct {
		tx *mock_transaction.MockTx
	}
	type args struct {
		ctx context.Context
	}
	type expected struct {
		panic bool
	}
	tests := []struct {
		args     args
		prepare  func(*fields, *args)
		fields   fields
		name     string
		expected expected
	}{
		{
			name: "success",
			prepare: func(f *fields, a *args) {
				f.tx.EXPECT().
					Rollback(a.ctx).
					Return(nil)
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "logs_rollback_error",
			prepare: func(f *fields, a *args) {
				f.tx.EXPECT().
					Rollback(a.ctx).
					Return(assert.AnError)
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields.tx = mock_transaction.NewMockTx(ctrl)
			tt.prepare(&tt.fields, &tt.args)

			panicked := false
			func() {
				defer func() {
					panicked = recover() != nil
				}()
				rollback(tt.args.ctx, tt.fields.tx)
			}()

			assert.Equal(t, tt.expected.panic, panicked)
		})
	}
}

func Test_commit(t *testing.T) {
	t.Parallel()

	type fields struct {
		tx *mock_transaction.MockTx
	}
	type args struct {
		ctx context.Context
	}
	type expected struct {
		err      error
		errDebug []string
	}
	tests := []struct {
		args     args
		prepare  func(*fields, *args)
		fields   fields
		name     string
		expected expected
	}{
		{
			name: "success",
			prepare: func(f *fields, a *args) {
				f.tx.EXPECT().
					Commit(a.ctx).
					Return(nil)
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "commit_error",
			prepare: func(f *fields, a *args) {
				f.tx.EXPECT().
					Commit(a.ctx).
					Return(assert.AnError)
			},
			args: args{
				ctx: context.Background(),
			},
			expected: expected{
				err:      assert.AnError,
				errDebug: []string{"tx.Commit"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields.tx = mock_transaction.NewMockTx(ctrl)
			tt.prepare(&tt.fields, &tt.args)

			err := commit(tt.args.ctx, tt.fields.tx)

			requireError(t, err, tt.expected.err, tt.expected.errDebug)
		})
	}
}

func Test_txToContext(t *testing.T) {
	t.Parallel()

	type fields struct {
		tx *mock_transaction.MockTx
	}
	type args struct {
		ctx context.Context
	}
	type expected struct {
		value pgx.Tx
	}
	tests := []struct {
		prepare  func(*fields, *args, *expected)
		fields   fields
		args     args
		expected expected
		name     string
	}{
		{
			name: "stores_tx",
			prepare: func(f *fields, _ *args, e *expected) {
				e.value = f.tx
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields.tx = mock_transaction.NewMockTx(ctrl)
			tt.prepare(&tt.fields, &tt.args, &tt.expected)

			ctx := txToContext(tt.args.ctx, tt.fields.tx)
			value := txFromContext(ctx)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func Test_txFromContext(t *testing.T) {
	t.Parallel()

	type fields struct {
		tx *mock_transaction.MockTx
	}
	type args struct {
		ctx context.Context
	}
	type expected struct {
		value pgx.Tx
	}
	tests := []struct {
		prepare  func(*fields, *args, *expected)
		fields   fields
		args     args
		expected expected
		name     string
	}{
		{
			name: "empty_context",
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "context_with_tx",
			prepare: func(f *fields, a *args, e *expected) {
				a.ctx = txToContext(context.Background(), f.tx)
				e.value = f.tx
			},
		},
		{
			name: "context_with_wrong_value",
			args: args{
				ctx: context.WithValue(context.Background(), txContextKey, "tx"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields.tx = mock_transaction.NewMockTx(ctrl)
			if tt.prepare != nil {
				tt.prepare(&tt.fields, &tt.args, &tt.expected)
			}

			value := txFromContext(tt.args.ctx)

			assert.Equal(t, tt.expected.value, value)
		})
	}
}

func scanInt64(value int64) func(...any) error {
	return func(dest ...any) error {
		if len(dest) == 0 {
			return assert.AnError
		}
		id, ok := dest[0].(*int64)
		if !ok {
			return assert.AnError
		}
		*id = value
		return nil
	}
}

func requireError(t *testing.T, actual error, expected error, expectedDebug []string) {
	t.Helper()

	if expected == nil {
		require.NoError(t, actual)
		return
	}

	require.Error(t, actual)
	assert.Equal(t, expected, errutil.E(actual).GetError())
	assert.Equal(t, expectedDebug, errutil.E(actual).GetDebug())
}
