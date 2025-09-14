package dto

type ReqEmailChange struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type ReqEmailCheckQuery struct {
	Email string `form:"email" binding:"required,email"`
}

type RespEmailCheckQuery struct {
	IsRegistered bool `json:"is_registered"`
}

type ReqEmailVerification struct {
	Token string `form:"token" binding:"required"`
}
