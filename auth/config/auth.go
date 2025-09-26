package config

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type authConfig struct {
	BCryptCost          int
	JWTIssuer           string
	JWTAudiences        []string
	JWTSecret           string
	JWTDuration         int
	SessionDuration     int
	ResetTokenDuration  int
	VerifyTokenDuration int
	ChangeTokenDuration int
	LockDuration        int
}

var authCfg *authConfig

func loadAuthConfig() {
	authCfg = &authConfig{
		BCryptCost:          getNumberEnvWithFallback("BCRYPT_COST", bcrypt.DefaultCost),
		JWTIssuer:           getEnv("JWT_ISSUER"),
		JWTAudiences:        strings.Split(getEnv("JWT_AUDIENCES"), ","),
		JWTSecret:           getEnv("JWT_SECRET"),
		JWTDuration:         getNumberEnv("JWT_DURATION"),
		SessionDuration:     getNumberEnv("SESSION_DURATION"),
		ResetTokenDuration:  getNumberEnv("RESET_TOKEN_DURATION"),
		VerifyTokenDuration: getNumberEnv("VERIFY_TOKEN_DURATION"),
		ChangeTokenDuration: getNumberEnv("EMAIL_CHANGE_TOKEN_DURATION"),
		LockDuration:        getNumberEnv("LOCK_DURATION"),
	}
}

func AuthGetBCryptCost() (cost int) {
	return authCfg.BCryptCost
}

func AuthGetJWTIssuer() (issuer string) {
	return authCfg.JWTIssuer
}

func AuthGetJWTAudiences() (audiences []string) {
	return authCfg.JWTAudiences
}

func AuthGetJWTSecret() (secret string) {
	return authCfg.JWTSecret
}

func AuthGetJWTDuration() (duration int) {
	return authCfg.JWTDuration
}

func AuthGetSessionDuration() (duration int) {
	return authCfg.SessionDuration
}

func AuthGetResetTokenDuration() (duration int) {
	return authCfg.ResetTokenDuration
}

func AuthGetVerifyTokenDuration() (duration int) {
	return authCfg.VerifyTokenDuration
}

func AuthGetChangeTokenDuration() (duration int) {
	return authCfg.ChangeTokenDuration
}

func AuthGetLockDuration() (duration int) {
	return authCfg.LockDuration
}
