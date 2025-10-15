package microsoft

import (
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"golang.org/x/oauth2"
)

var endpoints = oauth2.Endpoint{
	AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
	TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
}

func NewProvider(cfg *configs.OAuth) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Microsoft.ClientID,
		ClientSecret: cfg.Microsoft.Secret,
		RedirectURL:  cfg.Microsoft.RedirectURL,
		Scopes:       []string{"openid", "profile", "email", "offline_access"},
		Endpoint:     endpoints,
	}
}
