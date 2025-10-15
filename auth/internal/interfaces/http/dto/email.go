package dto

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type ConfirmEmailChangeRequest struct {
	Token string `form:"token" binding:"required"`
}

type QueryEmailRequest struct {
	Email string `form:"email" binding:"required,email"`
}

type ConfirmEmailChangeResponse struct {
	Token string       `json:"token"`
	Auth  AuthResponse `json:"auth"`
}

type QueryEmailResponse struct {
	IsRegistered bool `json:"is_registered"`
}
