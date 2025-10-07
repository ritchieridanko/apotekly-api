package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/database"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "repository.user"

type UserRepository interface {
	Create(ctx context.Context, authID int64, data *entities.CreateUser) (user *entities.User, err error)
	GetByAuthID(ctx context.Context, authID int64) (user *entities.User, err error)
	GetUserID(ctx context.Context, authID int64) (userID uuid.UUID, err error)
	Update(ctx context.Context, authID int64, data *entities.UpdateUser) (user *entities.User, err error)
	UpdateProfilePicture(ctx context.Context, authID int64, profilePicture string) (user *entities.User, err error)
	Exists(ctx context.Context, authID int64) (exists bool, err error)
}

type userRepository struct {
	database *database.Database
}

func NewUserRepository(database *database.Database) UserRepository {
	return &userRepository{database}
}

func (r *userRepository) Create(ctx context.Context, authID int64, data *entities.CreateUser) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO users (
			auth_id, user_id, name, bio, sex, birthdate, phone, profile_picture
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			user_id, name, bio, sex, birthdate, phone, profile_picture,
			created_at, updated_at
	`

	row := r.database.QueryRow(
		ctx, query,
		authID, data.ID, data.Name, data.Bio, data.Sex,
		data.Birthdate, data.Phone, data.ProfilePicture,
	)

	var user entities.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create user: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &user, nil
}

func (r *userRepository) GetByAuthID(ctx context.Context, authID int64) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "GetByAuthID")
	defer span.End()

	query := `
		SELECT
			user_id, name, bio, sex, birthdate, phone, profile_picture,
			created_at, updated_at
		FROM users
		WHERE auth_id = $1 AND deleted_at IS NULL
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var user entities.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch user: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeUserNotFound, ce.MsgUserNotFound, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &user, nil
}

func (r *userRepository) GetUserID(ctx context.Context, authID int64) (uuid.UUID, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "GetUserID")
	defer span.End()

	query := "SELECT user_id FROM users WHERE auth_id = $1 AND deleted_at IS NULL"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var userID uuid.UUID
	if err := row.Scan(&userID); err != nil {
		wErr := fmt.Errorf("failed to fetch user id: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return uuid.Nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return uuid.Nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return userID, nil
}

func (r *userRepository) Update(ctx context.Context, authID int64, data *entities.UpdateUser) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "Update")
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
		err := fmt.Errorf("failed to update user: %w", ce.ErrNoFieldsProvided)
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgNoFieldsToUpdate, err)
	}
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
			UPDATE users
			SET %s
			WHERE auth_id = $%d AND deleted_at IS NULL
			RETURNING
				user_id, name, bio, sex, birthdate, phone, profile_picture,
				created_at, updated_at
		`, strings.Join(setClauses, ", "), argPos,
	)
	args = append(args, authID)

	row := r.database.QueryRow(ctx, query, args...)

	var user entities.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to update user: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &user, nil
}

func (r *userRepository) UpdateProfilePicture(ctx context.Context, authID int64, profilePicture string) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "UpdateProfilePicture")
	defer span.End()

	query := `
		UPDATE users
		SET profile_picture = $1, updated_at = NOW()
		WHERE auth_id = $2 AND deleted_at IS NULL
		RETURNING
			user_id, name, bio, sex, birthdate, phone, profile_picture,
			created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, profilePicture, authID)

	var user entities.User
	err := row.Scan(
		&user.ID, &user.Name, &user.Bio, &user.Sex,
		&user.Birthdate, &user.Phone, &user.ProfilePicture,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to update profile picture: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &user, nil
}

func (r *userRepository) Exists(ctx context.Context, authID int64) (bool, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "Exists")
	defer span.End()

	query := "SELECT 1 FROM users WHERE auth_id = $1 AND deleted_at IS NULL"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		wErr := fmt.Errorf("failed to fetch user: %w", err)
		return false, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return true, nil
}
