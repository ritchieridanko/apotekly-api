package services

import (
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
)

type JWTService struct {
	secret    string
	issuer    string
	audiences []string
	duration  time.Duration
}

func NewJWTService(cfg *configs.Auth) *JWTService {
	return &JWTService{cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.Audiences, cfg.JWT.Duration}
}

func (s *JWTService) Create(authID int64, roleID int16, isVerified bool) (string, error) {
	now := time.Now().UTC()

	claim := entities.Claim{
		AuthID:     authID,
		RoleID:     roleID,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", authID),
			Audience:  jwt.ClaimStrings(s.audiences),
			IssuedAt:  &jwt.NumericDate{Time: now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(s.duration)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) Parse(tokenString string) (*entities.Claim, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&entities.Claim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.secret), nil
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

func (s *JWTService) ValidateAudience(appName string) bool {
	return slices.Contains(s.audiences, appName)
}
