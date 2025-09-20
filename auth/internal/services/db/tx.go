package db

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"go.opentelemetry.io/otel"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) (err error)) (err error)
}

type txManager struct {
	instance *sql.DB
}

func NewTxManager(instance *sql.DB) TxManager {
	return &txManager{instance}
}

func (m *txManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, span := otel.Tracer("service.db.tx_manager").Start(ctx, "WithTx")
	defer span.End()

	tx := getTxFromContext(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = m.instance.BeginTx(ctx, nil)
		if err != nil {
			return ce.NewError(span, ce.CodeDBTransaction, ce.MsgInternalServer, err)
		}
		ctx = setTxInContext(ctx, tx)
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
			return ce.NewError(span, ce.CodeDBTransaction, ce.MsgInternalServer, err)
		}
	}

	return nil
}
