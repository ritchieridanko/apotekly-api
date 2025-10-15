package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const sessionErrorTracer string = "repository.session"

type SessionRepository interface {
	Create(ctx context.Context, authID int64, data *entities.CreateSession) (err error)
	GetByToken(ctx context.Context, token string) (session *entities.Session, err error)
	RevokeActive(ctx context.Context, authID int64) (revokedSessionID int64, err error)
	RevokeByID(ctx context.Context, sessionID int64) (err error)
	RevokeByToken(ctx context.Context, token string) (err error)
}

type sessionRepository struct {
	database *database.Database
}

func NewSessionRepository(database *database.Database) SessionRepository {
	return &sessionRepository{database}
}

func (r *sessionRepository) Create(ctx context.Context, authID int64, data *entities.CreateSession) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO sessions (
			auth_id, parent_id, token, user_agent, ip_address, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	err := r.database.Execute(
		ctx, query,
		authID, data.ParentID, data.Token, data.UserAgent, data.IPAddress, data.ExpiresAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create session: %w", err)
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*entities.Session, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "GetByToken")
	defer span.End()

	query := `
		SELECT
			session_id, auth_id, parent_id, token, user_agent, ip_address,
			created_at, expires_at, revoked_at
		FROM sessions
		WHERE token = $1 AND revoked_at IS NULL
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, token)

	var session entities.Session
	err := row.Scan(
		&session.ID, &session.AuthID, &session.ParentID,
		&session.Token, &session.UserAgent, &session.IPAddress,
		&session.CreatedAt, &session.ExpiresAt, &session.RevokedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch session by token: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &session, nil
}

func (r *sessionRepository) RevokeActive(ctx context.Context, authID int64) (int64, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeActive")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE auth_id = $1 AND revoked_at IS NULL AND expires_at >= $2
		RETURNING session_id
	`

	row := r.database.QueryRow(ctx, query, authID, time.Now().UTC())

	var sessionID int64
	if err := row.Scan(&sessionID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		wErr := fmt.Errorf("failed to revoke active session: %w", err)
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return sessionID, nil
}

func (r *sessionRepository) RevokeByID(ctx context.Context, sessionID int64) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeByID")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE session_id = $1 AND revoked_at IS NULL
	`

	if err := r.database.Execute(ctx, query, sessionID); err != nil {
		wErr := fmt.Errorf("failed to revoke session by id: %w", err)
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}
	return nil
}

func (r *sessionRepository) RevokeByToken(ctx context.Context, token string) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeByToken")
	defer span.End()

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`

	if err := r.database.Execute(ctx, query, token); err != nil {
		wErr := fmt.Errorf("failed to revoke session by token: %w", err)
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}
	return nil
}
