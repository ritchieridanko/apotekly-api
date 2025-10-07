package database

import (
	"context"
	"database/sql"
)

type executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}

func (d *Database) getQueryExecutor(ctx context.Context) executor {
	if tx := txFromCtx(ctx); tx != nil {
		return tx
	}
	return d.db
}
