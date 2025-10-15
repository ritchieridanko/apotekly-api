package services

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CookieService struct {
	isSecure bool
	httpOnly bool
}

func NewCookieService(appEnv string, httpOnly bool) *CookieService {
	isSecure := strings.ToLower(strings.TrimSpace(appEnv)) == "production"
	return &CookieService{isSecure, httpOnly}
}

func (s *CookieService) Set(ctx *gin.Context, name, value string, duration time.Duration, path, domain string) {
	ctx.SetCookie(name, value, int(duration.Seconds()), path, domain, s.isSecure, s.httpOnly)
}

func (s *CookieService) Delete(ctx *gin.Context, name, path, domain string) {
	ctx.SetCookie(name, "", -1, path, domain, s.isSecure, s.httpOnly)
}
