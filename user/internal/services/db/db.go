package db

import (
	"context"
	"database/sql"
	"errors"
)

type Database interface {
	Execute(ctx context.Context, query string, args ...any) (err error)
	QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row)

	IsWithinTx(ctx context.Context) (isWithin bool)
}

type database struct {
	instance *sql.DB
}

func NewDatabase(instance *sql.DB) (db Database) {
	return &database{instance}
}

func (db *database) Execute(ctx context.Context, query string, args ...any) error {
	executor := getQueryExecutor(ctx, db.instance)
	result, err := executor.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("query execution affected no rows")
	}

	return nil
}

func (db *database) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	executor := getQueryExecutor(ctx, db.instance)
	return executor.QueryRowContext(ctx, query, args...)
}

func (db *database) IsWithinTx(ctx context.Context) bool {
	return getTxFromContext(ctx) != nil
}
