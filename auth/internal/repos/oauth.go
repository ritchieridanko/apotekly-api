package repos

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/db"
	"go.opentelemetry.io/otel"
)

const oAuthErrorTracer string = "repo.oauth"

type OAuthRepo interface {
	Create(ctx context.Context, authID int64, data *entities.OAuth) (oauthID int64, err error)
}

type oAuthRepo struct {
	database db.DBService
}

func NewOAuthRepo(database db.DBService) OAuthRepo {
	return &oAuthRepo{database}
}

func (r *oAuthRepo) Create(ctx context.Context, authID int64, data *entities.OAuth) (int64, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO oauth (auth_id, provider, provider_uid)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, authID, data.Provider, data.UID)

	var oauthID int64
	if err := row.Scan(&oauthID); err != nil {
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return oauthID, nil
}
