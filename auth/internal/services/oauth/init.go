package oauth

import "golang.org/x/oauth2"

func Initialize() (google, microsoft *oauth2.Config) {
	return initGoogleOAuth(), initMicrosoftOAuth()
}
