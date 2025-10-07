package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name      string                `form:"name" binding:"required,name"`
	Bio       *string               `form:"bio" binding:"omitempty,bio"`
	Sex       *string               `form:"sex" binding:"omitempty,sex"`
	Birthdate *time.Time            `form:"birthdate" time_format:"2006-01-02" binding:"omitempty,birthdate"`
	Phone     *string               `form:"phone" binding:"omitempty,phone"`
	Image     *multipart.FileHeader `form:"image"`
}

type UpdateUserRequest struct {
	Name      *string    `json:"name" binding:"omitempty,name"`
	Bio       *string    `json:"bio" binding:"omitempty,bio"`
	Sex       *string    `json:"sex" binding:"omitempty,sex"`
	Birthdate *time.Time `json:"birthdate" binding:"omitempty,birthdate"`
	Phone     *string    `json:"phone" binding:"omitempty,phone"`
}

type UserResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Bio            *string    `json:"bio"`
	Sex            *string    `json:"sex"`
	Birthdate      *time.Time `json:"birthdate"`
	Phone          *string    `json:"phone"`
	ProfilePicture *string    `json:"profile_picture"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateUserResponse struct {
	Created UserResponse `json:"created"`
}

type UpdateUserResponse struct {
	Updated UserResponse `json:"updated"`
}
