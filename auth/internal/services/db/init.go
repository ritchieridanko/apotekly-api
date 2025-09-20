package db

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
)

type DBService interface {
	Execute(ctx context.Context, query string, args ...any) (err error)
	QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row)

	IsWithinTx(ctx context.Context) (isWithin bool)
}

type dbService struct {
	instance *sql.DB
}

func NewService(instance *sql.DB) DBService {
	return &dbService{instance}
}

func (dbs *dbService) Execute(ctx context.Context, query string, args ...any) error {
	executor := getQueryExecutor(ctx, dbs.instance)
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

func (dbs *dbService) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	executor := getQueryExecutor(ctx, dbs.instance)
	return executor.QueryRowContext(ctx, query, args...)
}

func (dbs *dbService) IsWithinTx(ctx context.Context) bool {
	return getTxFromContext(ctx) != nil
}
