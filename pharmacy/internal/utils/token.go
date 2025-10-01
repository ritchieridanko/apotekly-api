package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/apotekly-api/pharmacy/config"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/entities"
)

func ParseJWTToken(tokenString string) (claim *entities.Claim, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&entities.Claim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.AuthGetJWTSecret()), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*entities.Claim)
	if !ok {
		return nil, errors.New("invalid jwt token claim")
	}

	return claim, nil
}

func IsInTokenAudience(audiences jwt.ClaimStrings) (isIn bool) {
	for _, audience := range audiences {
		if audience == config.AppGetName() {
			return true
		}
	}
	return false
}
