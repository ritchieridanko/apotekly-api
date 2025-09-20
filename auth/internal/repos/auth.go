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

const authErrorTracer string = "repo.auth"

type AuthRepo interface {
	Create(ctx context.Context, data *entities.NewAuth) (authID int64, err error)
	CreateByOAuth(ctx context.Context, data *entities.NewAuth) (authID int64, err error)
	GetByEmail(ctx context.Context, email string) (auth *entities.Auth, err error)
	GetByID(ctx context.Context, authID int64) (auth *entities.Auth, err error)
	GetForOAuth(ctx context.Context, email string) (exists bool, auth *entities.Auth, err error)
	UpdateEmail(ctx context.Context, authID int64, email string) (err error)
	UpdatePassword(ctx context.Context, authID int64, password string) (err error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err error)
	VerifyEmail(ctx context.Context, authID int64) (err error)
	LockAccount(ctx context.Context, authID int64, until time.Time) (err error)
}

type authRepo struct {
	database db.DBService
}

func NewAuthRepo(database db.DBService) AuthRepo {
	return &authRepo{database}
}

func (r *authRepo) Create(ctx context.Context, data *entities.NewAuth) (int64, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO auth (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, data.Email, data.Password, data.Role)

	var authID int64
	if err := row.Scan(&authID); err != nil {
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return authID, nil
}

func (r *authRepo) CreateByOAuth(ctx context.Context, data *entities.NewAuth) (int64, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "CreateByOAuth")
	defer span.End()

	query := `
		INSERT INTO auth (email, role)
		VALUES ($1, $2)
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, data.Email, data.Role)

	var authID int64
	if err := row.Scan(&authID); err != nil {
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return authID, nil
}

func (r *authRepo) GetByEmail(ctx context.Context, email string) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "GetByEmail")
	defer span.End()

	query := `
		SELECT id, email, password, is_verified, locked_until, role
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var auth entities.Auth
	if err := row.Scan(&auth.ID, &auth.Email, &auth.Password, &auth.IsVerified, &auth.LockedUntil, &auth.Role); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &auth, nil
}

func (r *authRepo) GetByID(ctx context.Context, authID int64) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "GetByID")
	defer span.End()

	query := `
		SELECT id, email, password, is_verified, locked_until, role
		FROM auth
		WHERE id = $1 AND deleted_at IS NULL
	`
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var auth entities.Auth
	if err := row.Scan(&auth.ID, &auth.Email, &auth.Password, &auth.IsVerified, &auth.LockedUntil, &auth.Role); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &auth, nil
}

func (r *authRepo) GetForOAuth(ctx context.Context, email string) (bool, *entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "GetForOAuth")
	defer span.End()

	query := `
		SELECT id, email, password, is_verified, locked_until, role
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var auth entities.Auth
	if err := row.Scan(&auth.ID, &auth.Email, &auth.Password, &auth.IsVerified, &auth.LockedUntil, &auth.Role); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil, nil
		}
		return false, nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return true, &auth, nil
}

func (r *authRepo) UpdateEmail(ctx context.Context, authID int64, email string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UpdateEmail")
	defer span.End()

	query := `
		UPDATE auth
		SET email = $1, email_changed_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`

	if err := r.database.Execute(ctx, query, email, authID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (r *authRepo) UpdatePassword(ctx context.Context, authID int64, password string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UpdatePassword")
	defer span.End()

	query := `
		UPDATE auth
		SET password = $1, password_changed_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`

	if err := r.database.Execute(ctx, query, password, authID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (r *authRepo) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "IsEmailRegistered")
	defer span.End()

	query := "SELECT 1 FROM auth WHERE email = $1 AND deleted_at IS NULL"
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return true, nil
}

func (r *authRepo) VerifyEmail(ctx context.Context, authID int64) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "VerifyEmail")
	defer span.End()

	query := `
		UPDATE auth
		SET is_verified = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_verified = FALSE
	`

	if err := r.database.Execute(ctx, query, authID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (r *authRepo) LockAccount(ctx context.Context, authID int64, until time.Time) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "LockAccount")
	defer span.End()

	query := `
		UPDATE auth
		SET locked_until = $1, updated_at = NOW()
		WHERE id = $2
	`

	if err := r.database.Execute(ctx, query, until, authID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}
