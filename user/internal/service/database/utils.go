package database

import (
	"context"
	"database/sql"
)

type ctxKeyTx struct{}

var key ctxKeyTx = ctxKeyTx{}

func txToCtx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, key, tx)
}

func txFromCtx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(key).(*sql.Tx); ok {
		return tx
	}
	return nil
}
