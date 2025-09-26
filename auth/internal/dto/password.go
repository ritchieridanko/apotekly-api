package dto

type ReqChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

type ReqForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type ReqResetPassword struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}
