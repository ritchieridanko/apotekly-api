package repos

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "repo.user"

type UserRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewUser) (user *entities.User, err error)
	GetByAuthID(ctx context.Context, authID int64) (user *entities.User, err error)
	GetUserID(ctx context.Context, authID int64) (userID uuid.UUID, err error)
	UpdateUser(ctx context.Context, authID int64, data *entities.UserChange) (user *entities.User, err error)
	UpdateProfilePicture(ctx context.Context, authID int64, profilePicture string) (err error)
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
		ctx, query,
		authID, data.UserID, data.Name, data.Bio, data.Sex, data.Birthdate, data.Phone, data.ProfilePicture,
	)

	var user entities.User
	err := row.Scan(&user.UserID, &user.Name, &user.Bio, &user.Sex, &user.Birthdate, &user.Phone, &user.ProfilePicture)
	if err != nil {
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &user, nil
}

func (r *userRepo) GetByAuthID(ctx context.Context, authID int64) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "GetByAuthID")
	defer span.End()

	query := `
		SELECT user_id, name, bio, sex, birthdate, phone, profile_picture
		FROM users
		WHERE auth_id = $1 AND deleted_at IS NULL
	`
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var user entities.User
	err := row.Scan(&user.UserID, &user.Name, &user.Bio, &user.Sex, &user.Birthdate, &user.Phone, &user.ProfilePicture)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeUserNotFound, ce.MsgUserNotFound, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &user, nil
}

func (r *userRepo) GetUserID(ctx context.Context, authID int64) (uuid.UUID, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "GetUserID")
	defer span.End()

	query := "SELECT user_id FROM users WHERE auth_id = $1 AND deleted_at IS NULL"
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var userID uuid.UUID
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return uuid.Nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return uuid.Nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return userID, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, authID int64, data *entities.UserChange) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "UpdateUser")
	defer span.End()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if data.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argPos))
		args = append(args, *data.Bio)
		argPos++
	}
	if data.Sex != nil {
		setClauses = append(setClauses, fmt.Sprintf("sex = $%d", argPos))
		args = append(args, *data.Sex)
		argPos++
	}
	if data.Birthdate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birthdate = $%d", argPos))
		args = append(args, *data.Birthdate)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgNoFieldsToUpdate, ce.ErrNoFieldsProvided)
	}
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
			UPDATE users
			SET %s
			WHERE auth_id = $%d AND deleted_at IS NULL
			RETURNING user_id, name, bio, sex, birthdate, phone, profile_picture
		`, strings.Join(setClauses, ", "), argPos,
	)
	args = append(args, authID)

	row := r.database.QueryRow(ctx, query, args...)

	var user entities.User
	err := row.Scan(
		&user.UserID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &user, nil
}

func (r *userRepo) UpdateProfilePicture(ctx context.Context, authID int64, profilePicture string) error {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "UpdateProfilePicture")
	defer span.End()

	query := `
		UPDATE users
		SET profile_picture = $1, updated_at = NOW()
		WHERE auth_id = $2 AND deleted_at IS NULL
	`

	if err := r.database.Execute(ctx, query, profilePicture, authID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
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
