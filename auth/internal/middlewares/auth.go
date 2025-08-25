package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const AuthErrorTracer = ce.AuthMiddlewareTracer

func Authenticate() gin.HandlerFunc {
	tracer := AuthErrorTracer + ": Authenticate()"

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, ce.ErrTokenNotFound)
			ctx.Error(err)
			ctx.Abort()
			return
		}

		authLength := 2
		authSlice := strings.Split(authHeader, " ")
		if len(authSlice) != authLength {
			err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, ce.ErrInvalidTokenFormat)
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if authType := strings.ToLower(authSlice[0]); authType != "bearer" {
			err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, ce.ErrInvalidTokenFormat)
			ctx.Error(err)
			ctx.Abort()
			return
		}

		token := authSlice[1]
		claim, err := utils.ParseJWTToken(token)
		if err != nil {
			err := ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, jwt.ErrTokenExpired)
			}
			if errors.Is(err, jwt.ErrTokenMalformed) {
				err = ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, jwt.ErrTokenMalformed)
			}

			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set("authID", claim.AuthID)
		ctx.Set("roleID", claim.RoleID)
		ctx.Set("isVerified", claim.IsVerified)
		ctx.Next()
	}
}
