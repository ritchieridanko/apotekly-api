package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

const SessionErrorTracer = ce.SessionRepoTracer

type SessionRepo interface {
	Create(ctx context.Context, data *entities.NewSession) (sessionID int64, err error)
	Reissue(ctx context.Context, data *entities.ReissueSession) (newSessionID int64, err error)
	GetByToken(ctx context.Context, token string) (session *entities.Session, err error)
	RevokeByID(ctx context.Context, sessionID int64) (err error)
	RevokeByToken(ctx context.Context, token string) (err error)
	HasActiveSession(ctx context.Context, authID int64) (hasAny bool, sessionID int64, err error)
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

	var sessionID int64
	if err := row.Scan(&sessionID); err != nil {
		return 0, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return sessionID, nil
}

func (r *sessionRepo) Reissue(ctx context.Context, data *entities.ReissueSession) (int64, error) {
	tracer := SessionErrorTracer + ": Reissue()"

	query := `
		INSERT INTO sessions (auth_id, parent_id, token, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, data.AuthID, data.ParentID, data.Token, data.UserAgent, data.IPAddress, data.ExpiresAt)

	var newSessionID int64
	if err := row.Scan(&newSessionID); err != nil {
		return 0, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return newSessionID, nil
}

func (r *sessionRepo) GetByToken(ctx context.Context, token string) (*entities.Session, error) {
	tracer := SessionErrorTracer + ": GetByToken()"

	query := `
		SELECT id, auth_id, parent_id, token, user_agent, ip_address, created_at, expires_at, revoked_at
		FROM sessions
		WHERE token = $1
	`

	if dbtx.IsInsideTx(ctx) {
		query += " FOR UPDATE"
	}

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, token)

	var session entities.Session
	err := row.Scan(
		&session.ID,
		&session.AuthID,
		&session.ParentID,
		&session.Token,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, err)
		}
		return nil, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return &session, nil
}

func (r *sessionRepo) RevokeByID(ctx context.Context, sessionID int64) error {
	tracer := SessionErrorTracer + ": RevokeByID()"

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	result, err := executor.ExecContext(ctx, query, sessionID)
	if err != nil {
		return ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}
	if rowsAffected == 0 {
		return ce.NewError(ce.ErrCodeDBNoChange, ce.ErrMsgInvalidCredentials, tracer, ce.ErrDBNoChange)
	}

	return nil
}

func (r *sessionRepo) RevokeByToken(ctx context.Context, token string) error {
	tracer := SessionErrorTracer + ": RevokeByToken()"

	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE token = $1 AND revoked_at IS NULL
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	result, err := executor.ExecContext(ctx, query, token)
	if err != nil {
		return ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}
	if rowsAffected == 0 {
		return ce.NewError(ce.ErrCodeDBNoChange, ce.ErrMsgInvalidCredentials, tracer, ce.ErrDBNoChange)
	}

	return nil
}

func (r *sessionRepo) HasActiveSession(ctx context.Context, authID int64) (bool, int64, error) {
	tracer := SessionErrorTracer + ": HasActiveSession()"

	query := `
		SELECT id, expires_at
		FROM sessions
		WHERE auth_id = $1 AND revoked_at IS NULL
	`

	if dbtx.IsInsideTx(ctx) {
		query += " FOR UPDATE"
	}

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, authID)

	var sessionID int64
	var expiresAt time.Time
	if err := row.Scan(&sessionID, &expiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, 0, nil
		}
		return false, 0, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return expiresAt.After(time.Now().UTC()), sessionID, nil
}
