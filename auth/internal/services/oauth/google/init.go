package google

import (
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewProvider(cfg *configs.OAuth) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.Secret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
