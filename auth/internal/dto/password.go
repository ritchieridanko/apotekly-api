package dto

type ReqPasswordChange struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

type ReqForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}
