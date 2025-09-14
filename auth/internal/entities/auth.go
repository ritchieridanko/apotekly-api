package entities

import "time"

type Auth struct {
	ID          int64
	Email       string
	Password    *string
	IsVerified  bool
	LockedUntil *time.Time
	Role        int16
}

type NewAuth struct {
	Email    string
	Password string
	Role     int16
}

type GetAuth struct {
	Email    string
	Password string
}
