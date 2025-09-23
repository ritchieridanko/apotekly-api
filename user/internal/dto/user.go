package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type ReqNewUser struct {
	Name      string                `form:"name" binding:"required"`
	Bio       *string               `form:"bio"`
	Sex       *int16                `form:"sex"`
	Birthdate *time.Time            `form:"birthdate" time_format:"2006-01-02"`
	Phone     *string               `form:"phone"`
	Image     *multipart.FileHeader `form:"image"`
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

type ReqUserUpdate struct {
	Name      *string    `json:"name,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	Sex       *int16     `json:"sex,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
}
