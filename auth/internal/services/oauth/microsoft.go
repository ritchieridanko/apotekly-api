package oauth

import (
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"golang.org/x/oauth2"
)

var microsoftEndpoint = oauth2.Endpoint{
	AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
	TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
}

func InitMicrosoftOAuth() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GetMicrosoftClientID(),
		ClientSecret: config.GetMicrosoftClientSecret(),
		RedirectURL:  config.GetMicrosoftRedirectURL(),
		Scopes:       []string{"openid", "profile", "email", "offline_access"},
		Endpoint:     microsoftEndpoint,
	}
}
