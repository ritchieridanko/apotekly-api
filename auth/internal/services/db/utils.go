package db

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
)

func setTxInContext(ctx context.Context, tx *sql.Tx) (ctxWithTx context.Context) {
	return context.WithValue(ctx, constants.CtxKeyTx, tx)
}

func getTxFromContext(ctx context.Context) (txFromCtx *sql.Tx) {
	if tx, ok := ctx.Value(constants.CtxKeyTx).(*sql.Tx); ok {
		return tx
	}
	return nil
}

type queryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}

func getQueryExecutor(ctx context.Context, instance *sql.DB) (executor queryExecutor) {
	if tx := getTxFromContext(ctx); tx != nil {
		return tx
	}
	return instance
}
