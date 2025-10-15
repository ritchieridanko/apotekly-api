package dto

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

type QueryTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type QueryTokenResponse struct {
	IsValid bool `json:"is_valid"`
}
