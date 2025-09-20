package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
)

func GenerateRandomToken() (token string) {
	return uuid.New().String()
}

func GenerateJWTToken(authID int64, roleID int16, isVerified bool) (jwtToken string, err error) {
	now := time.Now().UTC()
	jwtDuration := time.Duration(config.AuthGetJWTDuration()) * time.Minute

	claim := entities.Claim{
		AuthID:     authID,
		RoleID:     roleID,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.AuthGetJWTIssuer(),
			Subject:   fmt.Sprintf("%d", authID),
			Audience:  jwt.ClaimStrings(config.AuthGetJWTAudiences()),
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(jwtDuration)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	jwtToken, err = token.SignedString([]byte(config.AuthGetJWTSecret()))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

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

func GenerateURLWithTokenQuery(path, token string) (url string) {
	return config.ClientGetBaseURL() + path + "?token=" + token
}
