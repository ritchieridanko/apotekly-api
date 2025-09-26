package dto

type ReqQueryToken struct {
	Token string `json:"token" binding:"required"`
}

type RespQueryToken struct {
	IsValid bool `json:"is_valid"`
}
