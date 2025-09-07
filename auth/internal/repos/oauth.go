package repos

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

const OAuthErrorTracer = ce.OAuthRepoTracer

type OAuthRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewOAuth) (oauthID int64, err error)
}

type oAuthRepo struct {
	db *sql.DB
}

func NewOAuthRepo(db *sql.DB) OAuthRepo {
	return &oAuthRepo{db}
}

func (r *oAuthRepo) Create(ctx context.Context, authID int64, data *entities.NewOAuth) (int64, error) {
	tracer := OAuthErrorTracer + ": Create()"

	query := `
		INSERT INTO oauth (auth_id, provider, provider_uid)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, authID, data.Provider, data.UID)

	var oauthID int64
	if err := row.Scan(&oauthID); err != nil {
		return 0, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return oauthID, nil
}
