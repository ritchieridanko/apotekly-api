package database

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db}
}

func (d *Database) Execute(ctx context.Context, query string, args ...any) error {
	executor := d.getQueryExecutor(ctx)
	result, err := executor.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ce.ErrDBAffectNoRows
	}

	return nil
}

func (d *Database) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	executor := d.getQueryExecutor(ctx)
	return executor.QueryRowContext(ctx, query, args...)
}

func (d *Database) QueryAll(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	executor := d.getQueryExecutor(ctx)
	return executor.QueryContext(ctx, query, args...)
}

func (d *Database) InTx(ctx context.Context) bool {
	return txFromCtx(ctx) != nil
}
