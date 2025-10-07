package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{db}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, span := otel.Tracer("database.transactor").Start(ctx, "WithTx")
	defer span.End()

	tx := txFromCtx(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = t.db.BeginTx(ctx, nil)
		if err != nil {
			wErr := fmt.Errorf("failed to begin transaction: %w", err)
			return ce.NewError(span, ce.CodeDBTransaction, ce.MsgInternalServer, wErr)
		}
		ctx = txToCtx(ctx, tx)
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
			wErr := fmt.Errorf("failed to commit transaction: %w", err)
			return ce.NewError(span, ce.CodeDBTransaction, ce.MsgInternalServer, wErr)
		}
	}

	return nil
}
