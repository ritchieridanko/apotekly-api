package dbtx

import (
	"context"
	"database/sql"
)

type TxManager interface {
	ReturnError(ctx context.Context, fn func(ctx context.Context) (err error)) (err error)
	ReturnAnyAndError(ctx context.Context, fn func(ctx context.Context) (result any, err error)) (result any, err error)
}

type txManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db}
}

func (m *txManager) ReturnError(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := GetTxFromContext(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		ctx = WithTxContext(ctx, tx)
		isNewTx = true
	}

	if err := fn(ctx); err != nil {
		if isNewTx {
			_ = tx.Rollback()
		}
		return err
	}

	if isNewTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (m *txManager) ReturnAnyAndError(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	tx := GetTxFromContext(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		ctx = WithTxContext(ctx, tx)
		isNewTx = true
	}

	result, err := fn(ctx)
	if err != nil {
		if isNewTx {
			_ = tx.Rollback()
		}
		return nil, err
	}

	if isNewTx {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}

	return result, nil
}
