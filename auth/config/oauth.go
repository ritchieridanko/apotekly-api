package config

type oAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftRedirectURL  string
}

var oAuthCfg *oAuthConfig

func loadOAuthConfig() {
	oAuthCfg = &oAuthConfig{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL"),

		MicrosoftClientID:     getEnv("MICROSOFT_CLIENT_ID"),
		MicrosoftClientSecret: getEnv("MICROSOFT_CLIENT_SECRET"),
		MicrosoftRedirectURL:  getEnv("MICROSOFT_REDIRECT_URL"),
	}
}

func OAuthGoogleGetClientID() (id string) {
	return oAuthCfg.GoogleClientID
}

func OAuthGoogleGetClientSecret() (secret string) {
	return oAuthCfg.GoogleClientSecret
}

func OAuthGoogleGetRedirectURL() (url string) {
	return oAuthCfg.GoogleRedirectURL
}

func OAuthMicrosoftGetClientID() (id string) {
	return oAuthCfg.MicrosoftClientID
}

func OAuthMicrosoftGetClientSecret() (secret string) {
	return oAuthCfg.MicrosoftClientSecret
}

func OAuthMicrosoftGetRedirectURL() (url string) {
	return oAuthCfg.MicrosoftRedirectURL
}
