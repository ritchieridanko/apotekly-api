package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID
	Name           string
	Bio            *string
	Sex            *int16
	Birthdate      *time.Time
	Phone          *string
	ProfilePicture *string
}

type NewUser struct {
	UserID         uuid.UUID
	Name           string
	Bio            *string
	Sex            *int16
	Birthdate      *time.Time
	Phone          *string
	ProfilePicture *string
}
