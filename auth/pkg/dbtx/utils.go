package dbtx

import (
	"context"
	"database/sql"
)

type txKey struct{}

var txCtxKey = txKey{}

func WithTxContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}

func GetTxFromContext(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txCtxKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func IsInsideTx(ctx context.Context) bool {
	return GetTxFromContext(ctx) != nil
}
