package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

const AuthErrorTracer = ce.AuthRepoTracer

type AuthRepo interface {
	Create(ctx context.Context, data *entities.NewAuth) (authID int64, err error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err error)
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
