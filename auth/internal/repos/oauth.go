package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

const OAuthErrorTracer = ce.OAuthRepoTracer

type OAuthRepo interface {
	IsAuthRegistered(ctx context.Context, authID int64) (exists bool, err error)
}

type oAuthRepo struct {
	db *sql.DB
}

func NewOAuthRepo(db *sql.DB) OAuthRepo {
	return &oAuthRepo{db}
}

func (r *oAuthRepo) IsAuthRegistered(ctx context.Context, authID int64) (bool, error) {
	tracer := OAuthErrorTracer + ": IsAuthRegistered()"

	query := `
		SELECT 1
		FROM oauth
		WHERE auth_id = $1 AND deleted_at IS NULL
	`

	if dbtx.IsInsideTx(ctx) {
		query += " FOR UPDATE"
	}

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return true, nil
}
