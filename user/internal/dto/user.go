package dto

import (
	"time"

	"github.com/google/uuid"
)

type ReqNewUser struct {
	Name      string     `json:"name" binding:"required"`
	Bio       *string    `json:"bio,omitempty"`
	Sex       *int16     `json:"sex,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
}

type RespNewUser struct {
	UserID         uuid.UUID  `json:"user_id"`
	Name           string     `json:"name"`
	Bio            *string    `json:"bio"`
	Sex            *int16     `json:"sex"`
	Birthdate      *time.Time `json:"birthdate"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
}
