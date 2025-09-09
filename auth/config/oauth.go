package config

type oAuthConfig struct {
	googleClientID     string
	googleClientSecret string
	googleRedirectURL  string

	microsoftClientID     string
	microsoftClientSecret string
	microsoftRedirectURL  string
}

var oAuthCfg *oAuthConfig

func LoadOAuthConfig() {
	oAuthCfg = &oAuthConfig{
		googleClientID:     GetEnv("GOOGLE_CLIENT_ID"),
		googleClientSecret: GetEnv("GOOGLE_CLIENT_SECRET"),
		googleRedirectURL:  GetEnv("GOOGLE_REDIRECT_URL"),

		microsoftClientID:     GetEnv("MICROSOFT_CLIENT_ID"),
		microsoftClientSecret: GetEnv("MICROSOFT_CLIENT_SECRET"),
		microsoftRedirectURL:  GetEnv("MICROSOFT_REDIRECT_URL"),
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

func GetMicrosoftClientID() (clientID string) {
	return oAuthCfg.microsoftClientID
}

func GetMicrosoftClientSecret() (secret string) {
	return oAuthCfg.microsoftClientSecret
}

func GetMicrosoftRedirectURL() (url string) {
	return oAuthCfg.microsoftRedirectURL
}
