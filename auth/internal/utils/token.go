package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
)

func GenerateRandomToken() (token string) {
	return uuid.New().String()
}

func GenerateJWTToken(authID int64, roleID int16, isVerified bool) (tokenString string, err error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(config.GetJWTDuration()) * time.Minute)

	claim := entities.Claim{
		AuthID:     authID,
		RoleID:     roleID,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  config.GetJWTIssuer(),
			Subject: fmt.Sprintf("%d", authID),
			Audience: jwt.ClaimStrings{
				constants.AudienceAuthService,
				constants.AudienceCartService,
				constants.AudienceUserService,
			},
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err = token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
