package config

type authConfig struct {
	jwtSecret string
}

var authCfg *authConfig

func LoadAuthConfig() {
	authCfg = &authConfig{
		jwtSecret: GetEnv("JWT_SECRET"),
	}
}

func GetJWTSecret() (secret string) {
	return authCfg.jwtSecret
}
