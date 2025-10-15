package dto

type AuthenticateRequest struct {
	Code string `form:"code" binding:"required"`
}

type ExchangeCodeRequest struct {
	Code string `form:"code" binding:"required"`
}

type ExchangeCodeResponse struct {
	Token string       `json:"token"`
	Auth  AuthResponse `json:"auth"`
}

type GoogleUser struct {
	ID         string `json:"id" binding:"required"`
	Email      string `json:"email" binding:"required"`
	IsVerified bool   `json:"verified_email" binding:"required"`
}

type MicrosoftUser struct {
	ID                string `json:"id"`
	UserPrincipalName string `json:"userPrincipalName"`
	Mail              string `json:"mail"`
}
