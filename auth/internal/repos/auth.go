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

const AuthErrorTracer = ce.AuthRepoTracer

type AuthRepo interface {
	Create(ctx context.Context, data *entities.NewAuth) (authID int64, err error)
	GetByEmail(ctx context.Context, email string) (auth *entities.Auth, err error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err error)
	LockAccount(ctx context.Context, authID int64, until time.Time) (err error)
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) AuthRepo {
	return &authRepo{db}
}

func (r *authRepo) Create(ctx context.Context, data *entities.NewAuth) (int64, error) {
	tracer := AuthErrorTracer + ": Create()"

	query := `
		INSERT INTO auth (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, data.Email, data.Password, data.Role)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return id, nil
}

func (r *authRepo) GetByEmail(ctx context.Context, email string) (*entities.Auth, error) {
	tracer := AuthErrorTracer + ": GetByEmail()"

	query := `
		SELECT id, email, password, is_verified, locked_until, role
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`

	if dbtx.IsInsideTx(ctx) {
		query += " FOR UPDATE"
	}

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, email)

	var auth entities.Auth
	if err := row.Scan(&auth.ID, &auth.Email, &auth.Password, &auth.IsVerified, &auth.LockedUntil, &auth.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, err)
		}
		return nil, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return &auth, nil
}

func (r *authRepo) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	tracer := AuthErrorTracer + ": IsEmailRegistered()"

	query := `
		SELECT 1
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`

	if dbtx.IsInsideTx(ctx) {
		query += " FOR UPDATE"
	}

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, email)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return true, nil
}

func (r *authRepo) LockAccount(ctx context.Context, authID int64, until time.Time) error {
	tracer := AuthErrorTracer + ": LockAccount()"

	query := `
		UPDATE auth
		SET locked_until = $1, updated_at = NOW()
		WHERE id = $2
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	result, err := executor.ExecContext(ctx, query, until, authID)
	if err != nil {
		return ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}
	if rowsAffected == 0 {
		return ce.NewError(ce.ErrCodeDBNoChange, ce.ErrMsgInternalServer, tracer, ce.ErrDBNoChange)
	}

	return nil
}
