package dto

type ReqOAuthByGoogle struct {
	Code string `form:"code" binding:"required"`
}

type RespOAuthByGoogle struct {
	ID         string `json:"id" binding:"required"`
	Email      string `json:"email" binding:"required"`
	IsVerified bool   `json:"verified_email" binding:"required"`
}
