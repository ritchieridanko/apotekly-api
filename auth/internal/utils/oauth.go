package utils

import (
	"context"
	"encoding/json"

	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"golang.org/x/oauth2"
)

func GetUserInfoFromGoogle(ctx context.Context, token *oauth2.Token, config *oauth2.Config) (userInfo *dto.RespOAuthByGoogle, err error) {
	client := config.Client(ctx, token)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var info dto.RespOAuthByGoogle
	if err := json.NewDecoder(response.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
