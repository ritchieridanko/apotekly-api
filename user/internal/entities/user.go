package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Bio            *string
	Sex            *string
	Birthdate      *time.Time
	Phone          *string
	ProfilePicture *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateUser struct {
	ID             uuid.UUID
	Name           string
	Bio            *string
	Sex            *string
	Birthdate      *time.Time
	Phone          *string
	ProfilePicture *string
}

type UpdateUser struct {
	Name      *string
	Bio       *string
	Sex       *string
	Birthdate *time.Time
	Phone     *string
}
