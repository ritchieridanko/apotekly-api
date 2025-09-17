package config

type authConfig struct {
	JWTSecret string
}

var authCfg *authConfig

func loadAuthConfig() {
	authCfg = &authConfig{
		JWTSecret: getEnv("JWT_SECRET"),
	}
}

func AuthGetJWTSecret() (secret string) {
	return authCfg.JWTSecret
}
