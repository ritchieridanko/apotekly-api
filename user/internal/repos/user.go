package repos

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "repo.user"

type UserRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewUser) (user *entities.User, err error)
	HasUser(ctx context.Context, authID int64) (exists bool, err error)
}

type userRepo struct {
	database db.DBService
}

func NewUserRepo(database db.DBService) UserRepo {
	return &userRepo{database}
}

func (r *userRepo) Create(ctx context.Context, authID int64, data *entities.NewUser) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO users (auth_id, user_id, name, bio, sex, birthdate, phone, profile_picture)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING user_id, name, bio, sex, birthdate, phone, profile_picture
	`

	row := r.database.QueryRow(
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

	var user entities.User
	err := row.Scan(
		&user.UserID,
		&user.Name,
		&user.Bio,
		&user.Sex,
		&user.Birthdate,
		&user.Phone,
		&user.ProfilePicture,
	)
	if err != nil {
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &user, nil
}

func (r *userRepo) HasUser(ctx context.Context, authID int64) (bool, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "HasUser")
	defer span.End()

	query := "SELECT 1 FROM users WHERE auth_id = $1 AND deleted_at IS NULL"
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return true, nil
}
