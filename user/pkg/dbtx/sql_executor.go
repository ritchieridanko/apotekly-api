package dbtx

import (
	"context"
	"database/sql"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}

func GetSQLExecutor(ctx context.Context, db *sql.DB) (executor SQLExecutor) {
	if tx := GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return db
}
