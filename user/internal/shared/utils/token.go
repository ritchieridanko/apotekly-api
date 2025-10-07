package utils

import (
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
)

func JWTTokenParse(tokenString, jwtSecret string) (*entities.Claim, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&entities.Claim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*entities.Claim)
	if !ok {
		return nil, ce.ErrInvalidTokenClaim
	}

	return claim, nil
}

func JWTTokenValidateAudience(audiences jwt.ClaimStrings, appName string) bool {
	return slices.Contains(audiences, appName)
}
