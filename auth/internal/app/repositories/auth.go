package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const authErrorTracer string = "repository.auth"

type AuthRepository interface {
	Create(ctx context.Context, data *entities.CreateAuth) (createdAuth *entities.Auth, err error)
	GetByEmail(ctx context.Context, email string) (auth *entities.Auth, err error)
	GetByID(ctx context.Context, authID int64) (auth *entities.Auth, err error)
	GetForOAuth(ctx context.Context, email string) (exists bool, auth *entities.Auth, err error)
	UpdateEmail(ctx context.Context, authID int64, email string) (updatedAuth *entities.Auth, err error)
	UpdatePassword(ctx context.Context, authID int64, password string) (err error)
	SetVerified(ctx context.Context, authID int64) (verifiedAuth *entities.Auth, err error)
	Exists(ctx context.Context, email string) (exists bool, err error)
}

type authRepository struct {
	database *database.Database
}

func NewAuthRepository(database *database.Database) AuthRepository {
	return &authRepository{database}
}

func (r *authRepository) Create(ctx context.Context, data *entities.CreateAuth) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO auth (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING auth_id, email, role, is_verified, created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, data.Email, data.Password, data.RoleID)

	var auth entities.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.RoleID, &auth.IsVerified,
		&auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create auth: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &auth, nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "GetByEmail")
	defer span.End()

	query := `
		SELECT
			auth_id, email, password, role, is_verified,
			created_at, updated_at
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var auth entities.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.Password, &auth.RoleID,
		&auth.IsVerified, &auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth by email: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &auth, nil
}

func (r *authRepository) GetByID(ctx context.Context, authID int64) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "GetByID")
	defer span.End()

	query := `
		SELECT
			auth_id, email, password, role, is_verified,
			created_at, updated_at
		FROM auth
		WHERE auth_id = $1 AND deleted_at IS NULL
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var auth entities.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.Password, &auth.RoleID,
		&auth.IsVerified, &auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth by id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &auth, nil
}

func (r *authRepository) GetForOAuth(ctx context.Context, email string) (bool, *entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "GetForOAuth")
	defer span.End()

	query := `
		SELECT
			auth_id, email, password, role, is_verified,
			created_at, updated_at
		FROM auth
		WHERE email = $1 AND deleted_at IS NULL
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var auth entities.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.Password, &auth.RoleID,
		&auth.IsVerified, &auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil, nil
		}
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return false, nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return true, &auth, nil
}

func (r *authRepository) UpdateEmail(ctx context.Context, authID int64, email string) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UpdateEmail")
	defer span.End()

	query := `
		UPDATE auth
		SET email = $1, email_changed_at = NOW(), updated_at = NOW()
		WHERE auth_id = $2 AND deleted_at IS NULL
		RETURNING auth_id, email, role, is_verified, created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, email, authID)

	var auth entities.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.RoleID, &auth.IsVerified,
		&auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to update email: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &auth, nil
}

func (r *authRepository) UpdatePassword(ctx context.Context, authID int64, password string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UpdatePassword")
	defer span.End()

	query := `
		UPDATE auth
		SET password = $1, password_changed_at = NOW(), updated_at = NOW()
		WHERE auth_id = $2 AND deleted_at IS NULL
	`

	if err := r.database.Execute(ctx, query, password, authID); err != nil {
		wErr := fmt.Errorf("failed to update password: %w", err)
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}
	return nil
}

func (r *authRepository) SetVerified(ctx context.Context, authID int64) (*entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "SetVerified")
	defer span.End()

	query := `
		UPDATE auth
		SET is_verified = TRUE, updated_at = NOW()
		WHERE auth_id = $1 AND is_verified = FALSE AND deleted_at IS NULL
		RETURNING auth_id, email, role, is_verified, created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, authID)

	var auth entities.Auth
	err := row.Scan(
		&auth.ID, &auth.Email, &auth.RoleID, &auth.IsVerified,
		&auth.CreatedAt, &auth.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to set auth verified: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &auth, nil
}

func (r *authRepository) Exists(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Exists")
	defer span.End()

	query := "SELECT 1 FROM auth WHERE email = $1 AND deleted_at IS NULL"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, email)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return false, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return true, nil
}
