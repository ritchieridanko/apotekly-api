package dbtx

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/user/pkg/ce"
)

const DBTXErrorTracer = ce.DBTXTracer

type TxManager interface {
	ReturnError(ctx context.Context, fn func(ctx context.Context) (err error)) (err error)
}

type txManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db}
}

func (m *txManager) ReturnError(ctx context.Context, fn func(ctx context.Context) error) error {
	tracer := DBTXErrorTracer + ": ReturnError()"

	tx := GetTxFromContext(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = m.db.BeginTx(ctx, nil)
		if err != nil {
			return ce.NewError(ce.ErrCodeDBTX, ce.ErrMsgInternalServer, tracer, err)
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
			return ce.NewError(ce.ErrCodeDBTX, ce.ErrMsgInternalServer, tracer, err)
		}
	}

	return nil
}
