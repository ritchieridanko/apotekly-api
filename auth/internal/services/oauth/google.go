package oauth

import (
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func initGoogleOAuth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.OAuthGoogleGetClientID(),
		ClientSecret: config.OAuthGoogleGetClientSecret(),
		RedirectURL:  config.OAuthGoogleGetRedirectURL(),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
