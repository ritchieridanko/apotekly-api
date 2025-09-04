package dto

type ReqTokenCheckQuery struct {
	Token string `json:"token" binding:"required"`
}

type RespTokenCheckQuery struct {
	IsValid bool `json:"is_valid"`
}
