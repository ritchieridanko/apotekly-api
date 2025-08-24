package repos

import (
	"context"
	"database/sql"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

const SessionErrorTracer = ce.SessionRepoTracer

type SessionRepo interface {
	Create(ctx context.Context, data *entities.NewSession) (sessionID int64, err error)
}

type sessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) SessionRepo {
	return &sessionRepo{db}
}

func (r *sessionRepo) Create(ctx context.Context, data *entities.NewSession) (int64, error) {
	tracer := SessionErrorTracer + ": Create()"

	query := `
		INSERT INTO sessions (auth_id, token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, data.AuthID, data.Token, data.UserAgent, data.IPAddress, data.ExpiresAt)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return id, nil
}
