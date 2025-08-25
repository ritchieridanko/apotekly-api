package dto

type ReqEmailChange struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}
