package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
)

func SetSessionCookie(ctx *gin.Context, token string) {
	second := 60
	duration := config.GetSessionDuration() * second
	isSecure := strings.ToLower(config.GetServerProtocol()) == "https"
	host := config.GetServerHost()
	ctx.SetCookie(constants.CookieKeySessionToken, token, duration, "/", host, isSecure, true)
}
