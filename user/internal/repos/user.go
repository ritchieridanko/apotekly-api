package repos

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/user/pkg/dbtx"
)

const UserErrorTracer = ce.UserRepoTracer

type UserRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewUser) (err error)
	HasUser(ctx context.Context, authID int64) (exists bool, err error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db}
}

func (r *userRepo) Create(ctx context.Context, authID int64, data *entities.NewUser) error {
	tracer := UserErrorTracer + ": Create()"

	query := `
		INSERT INTO users (auth_id, user_id, name, bio, sex, birthdate, phone, profile_picture)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	result, err := executor.ExecContext(
		ctx,
		query,
		authID,
		data.UserID,
		data.Name,
		data.Bio,
		data.Sex,
		data.Birthdate,
		data.Phone,
		data.ProfilePicture,
	)
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

func (r *userRepo) HasUser(ctx context.Context, authID int64) (bool, error) {
	tracer := UserErrorTracer + ": HasUser()"

	query := "SELECT 1 FROM users WHERE auth_id = $1"

	if dbtx.IsInsideTx(ctx) {
		query += " FOR UPDATE"
	}

	executor := dbtx.GetSQLExecutor(ctx, r.db)
	row := executor.QueryRowContext(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, ce.NewError(ce.ErrCodeDBQuery, ce.ErrMsgInternalServer, tracer, err)
	}

	return true, nil
}
