package dto

type ReqChangeEmail struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type ReqQueryEmail struct {
	Email string `form:"email" binding:"required,email"`
}

type RespQueryEmail struct {
	IsRegistered bool `json:"is_registered"`
}

type ReqVerifyEmail struct {
	Token string `form:"token" binding:"required"`
}
