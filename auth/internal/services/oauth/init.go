package oauth

import (
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/oauth/google"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/oauth/microsoft"
	"golang.org/x/oauth2"
)

type OAuth struct {
	google    *oauth2.Config
	microsoft *oauth2.Config
}

func Initialize(cfg *configs.OAuth) *OAuth {
	g := google.NewProvider(cfg)
	m := microsoft.NewProvider(cfg)
	return &OAuth{google: g, microsoft: m}
}

func (o *OAuth) Google() *oauth2.Config {
	return o.google
}

func (o *OAuth) Microsoft() *oauth2.Config {
	return o.microsoft
}
