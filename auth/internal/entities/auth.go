package entities

import "time"

type Auth struct {
	ID                int64
	Email             string
	Password          *string
	RoleID            int16
	IsVerified        bool
	EmailChangedAt    *time.Time
	PasswordChangedAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateAuth struct {
	Email    string
	Password *string
	RoleID   int16
}

type GetAuth struct {
	Email    string
	Password string
}
