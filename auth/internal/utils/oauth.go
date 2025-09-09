package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"golang.org/x/oauth2"
)

func GetUserFromGoogle(ctx context.Context, token *oauth2.Token, config *oauth2.Config) (user *dto.GoogleUser, err error) {
	client := config.Client(ctx, token)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: %s", response.Status)
	}

	var googleUser dto.GoogleUser
	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func GetUserFromMicrosoft(ctx context.Context, accessToken string) (user *dto.MicrosoftUser, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: %s", response.Status)
	}

	var microsoftUser dto.MicrosoftUser
	if err := json.NewDecoder(response.Body).Decode(&microsoftUser); err != nil {
		return nil, err
	}

	return &microsoftUser, nil
}
