package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/apotekly-api/user/config"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
)

func ParseJWTToken(tokenString string) (claim *entities.Claim, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&entities.Claim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*entities.Claim)
	if !ok {
		return nil, fmt.Errorf("invalid jwt token type")
	}

	return claim, nil
}

func IsAudienceValid(audiences jwt.ClaimStrings) (isValid bool) {
	for _, audience := range audiences {
		if audience == config.GetAppName() {
			return true
		}
	}
	return false
}
