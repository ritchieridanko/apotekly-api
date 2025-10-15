package repositories

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const oAuthErrorTracer string = "repository.oauth"

type OAuthRepository interface {
	Create(ctx context.Context, authID int64, data *entities.OAuth) (err error)
}

type oAuthRepository struct {
	database *database.Database
}

func NewOAuthRepository(database *database.Database) OAuthRepository {
	return &oAuthRepository{database}
}

func (r *oAuthRepository) Create(ctx context.Context, authID int64, data *entities.OAuth) error {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO oauth (auth_id, provider, provider_uid)
		VALUES ($1, $2, $3)
	`

	if err := r.database.Execute(ctx, query, authID, data.Provider, data.UID); err != nil {
		wErr := fmt.Errorf("failed to create oauth: %w", err)
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}
	return nil
}
