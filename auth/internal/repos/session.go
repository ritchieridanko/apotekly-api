package repos

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/db"
	"go.opentelemetry.io/otel"
)

const sessionErrorTracer string = "repo.session"

type SessionRepo interface {
	Create(ctx context.Context, data *entities.NewSession) (sessionID int64, err error)
	Reissue(ctx context.Context, data *entities.SessionReissue) (newSessionID int64, err error)
	GetByToken(ctx context.Context, token string) (session *entities.Session, err error)
	RevokeByID(ctx context.Context, sessionID int64) (err error)
	RevokeByToken(ctx context.Context, token string) (err error)
	RevokeActive(ctx context.Context, authID int64) (revokedSessionID int64, err error)
}

type sessionRepo struct {
	database db.DBService
}

func NewSessionRepo(database db.DBService) SessionRepo {
	return &sessionRepo{database}
}

func (r *sessionRepo) Create(ctx context.Context, data *entities.NewSession) (int64, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO sessions (auth_id, token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, data.AuthID, data.Token, data.UserAgent, data.IPAddress, data.ExpiresAt)

	var sessionID int64
	if err := row.Scan(&sessionID); err != nil {
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return sessionID, nil
}

func (r *sessionRepo) Reissue(ctx context.Context, data *entities.SessionReissue) (int64, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "Reissue")
	defer span.End()

	query := `
		INSERT INTO sessions (auth_id, parent_id, token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, data.AuthID, data.ParentID, data.Token, data.UserAgent, data.IPAddress, data.ExpiresAt)

	var newSessionID int64
	if err := row.Scan(&newSessionID); err != nil {
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return newSessionID, nil
}

func (r *sessionRepo) GetByToken(ctx context.Context, token string) (*entities.Session, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "GetByToken")
	defer span.End()

	query := `
		SELECT id, auth_id, parent_id, token, user_agent, ip_address, created_at, expires_at, revoked_at
		FROM sessions
		WHERE token = $1
	`
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, token)

	var session entities.Session
	err := row.Scan(
		&session.ID, &session.AuthID, &session.ParentID, &session.Token, &session.UserAgent,
		&session.IPAddress, &session.CreatedAt, &session.ExpiresAt, &session.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &session, nil
}

func (r *sessionRepo) RevokeByID(ctx context.Context, sessionID int64) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeByID")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`

	if err := r.database.Execute(ctx, query, sessionID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (r *sessionRepo) RevokeByToken(ctx context.Context, token string) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeByToken")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`

	if err := r.database.Execute(ctx, query, token); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (r *sessionRepo) RevokeActive(ctx context.Context, authID int64) (int64, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeActive")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE auth_id = $1 AND revoked_at IS NULL AND expires_at >= $2
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, authID, time.Now().UTC())

	var revokedSessionID int64
	if err := row.Scan(&revokedSessionID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return revokedSessionID, nil
}
