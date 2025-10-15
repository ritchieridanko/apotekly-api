package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyAccountRequest struct {
	Token string `form:"token" binding:"required"`
}

type AuthResponse struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	RoleID     int16     `json:"role_id"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RegisterResponse struct {
	Token string       `json:"token"`
	Auth  AuthResponse `json:"auth"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	Auth  AuthResponse `json:"auth"`
}

type VerifyAccountResponse struct {
	Token string       `json:"token"`
	Auth  AuthResponse `json:"auth"`
}

type RefreshSessionResponse struct {
	Token string `json:"token"`
}
