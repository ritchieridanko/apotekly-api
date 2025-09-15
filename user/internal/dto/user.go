package dto

import "time"

type ReqNewUser struct {
	Name      string     `json:"name" binding:"required,min=3,max=50"`
	Bio       *string    `json:"bio" binding:"bio"`
	Sex       *int16     `json:"sex" binding:"sex"`
	Birthdate *time.Time `json:"birthdate" binding:"birthdate"`
	Phone     *string    `json:"phone" binding:"phone"`
}
