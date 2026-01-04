package transaction

import (
	"chattery/internal/client/postgres"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

type Manager struct {
	conn   *pgx.Conn
	client *postgres.Queries
}

func NewManager(conn *pgx.Conn) *Manager {
	return &Manager{
		conn:   conn,
		client: postgres.New(conn),
	}
}

func (m *Manager) Query(ctx context.Context) postgres.Querier {
	if tx := txFromContext(ctx); tx != nil {
		return m.client.WithTx(tx)
	}
	return m.client
}

func (m *Manager) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	if tx := txFromContext(ctx); tx != nil {
		return fn(ctx)
	}

	tx, err := m.conn.Begin(ctx)
	if err != nil {
		return errors.E(err).Debug("m.conn.Begin")
	}

	if err = fn(txToContext(ctx, tx)); err != nil {
		rollback(ctx, tx)
		return err
	}
	return commit(ctx, tx)
}

func rollback(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil {
		logger.ErrorCtx(ctx, err, "tx.Rollback")
	}
}

func commit(ctx context.Context, tx pgx.Tx) error {
	if err := tx.Commit(ctx); err != nil {
		return errors.E(err).Debug("tx.Commit")
	}
	return nil
}

type txContextKeyType struct{}

var txContextKey = txContextKeyType{}

func txToContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

func txFromContext(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txContextKey).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}
