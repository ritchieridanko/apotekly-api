package config

type oAuthConfig struct {
	googleClientID     string
	googleClientSecret string
	googleRedirectURL  string
}

var oAuthCfg *oAuthConfig

func LoadOAuthConfig() {
	oAuthCfg = &oAuthConfig{
		googleClientID:     GetEnv("GOOGLE_CLIENT_ID"),
		googleClientSecret: GetEnv("GOOGLE_CLIENT_SECRET"),
		googleRedirectURL:  GetEnv("GOOGLE_REDIRECT_URL"),
	}
}

func GetGoogleClientID() (clientID string) {
	return oAuthCfg.googleClientID
}

func GetGoogleClientSecret() (secret string) {
	return oAuthCfg.googleClientSecret
}

func GetGoogleRedirectURL() (url string) {
	return oAuthCfg.googleRedirectURL
}
