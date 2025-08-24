package entities

import "github.com/golang-jwt/jwt/v5"

type Claim struct {
	AuthID     int64
	RoleID     int16
	IsVerified bool
	jwt.RegisteredClaims
}
