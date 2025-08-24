package config

import (
	"golang.org/x/crypto/bcrypt"
)

type authConfig struct {
	bcryptCost                     int
	jwtIssuer                      string
	jwtSecret                      string
	jwtDuration                    int
	sessionDuration                int
	passwordResetTokenDuration     int
	emailVerificationTokenDuration int
	authLockDuration               int
}

var authCfg *authConfig

func LoadAuthConfig() {
	authCfg = &authConfig{
		bcryptCost:                     GetNumberEnvWithFallback("BCRYPT_COST", bcrypt.DefaultCost),
		jwtIssuer:                      GetEnv("JWT_ISSUER"),
		jwtSecret:                      GetEnv("JWT_SECRET"),
		jwtDuration:                    GetNumberEnv("JWT_DURATION"),
		sessionDuration:                GetNumberEnv("SESSION_DURATION"),
		passwordResetTokenDuration:     GetNumberEnv("PASSWORD_RESET_TOKEN_DURATION"),
		emailVerificationTokenDuration: GetNumberEnv("EMAIL_VERIFICATION_TOKEN_DURATION"),
		authLockDuration:               GetNumberEnv("AUTH_LOCK_DURATION"),
	}
}

func GetBCryptCost() (cost int) {
	return authCfg.bcryptCost
}

func GetJWTIssuer() (issuer string) {
	return authCfg.jwtIssuer
}

func GetJWTSecret() (secret string) {
	return authCfg.jwtSecret
}

func GetJWTDuration() (duration int) {
	return authCfg.jwtDuration
}

func GetSessionDuration() (duration int) {
	return authCfg.sessionDuration
}

func GetPasswordResetTokenDuration() (duration int) {
	return authCfg.passwordResetTokenDuration
}

func GetEmailVerificationTokenDuration() (duration int) {
	return authCfg.emailVerificationTokenDuration
}

func GetAuthLockDuration() (duration int) {
	return authCfg.authLockDuration
}
