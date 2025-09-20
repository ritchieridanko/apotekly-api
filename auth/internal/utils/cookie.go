package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
)

func CookieSetSession(ctx *gin.Context, token string) {
	second := 60
	duration := config.AuthGetSessionDuration() * second
	isSecure := strings.ToLower(config.ServerGetProtocol()) == "https"
	host := config.ServerGetHost()
	ctx.SetCookie(constants.CookieKeySessionToken, token, duration, "/", host, isSecure, true)
}

func CookieDelSession(ctx *gin.Context) {
	isSecure := strings.ToLower(config.ServerGetProtocol()) == "https"
	host := config.ServerGetHost()
	ctx.SetCookie(constants.CookieKeySessionToken, "", -1, "/", host, isSecure, true)
}
