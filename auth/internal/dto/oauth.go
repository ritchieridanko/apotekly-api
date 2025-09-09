package dto

type ReqOAuth struct {
	Code string `form:"code" binding:"required"`
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
