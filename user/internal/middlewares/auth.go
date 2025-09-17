package middlewares

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/constants"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"go.opentelemetry.io/otel"
)

const AuthErrorTracer string = "auth.middleware"

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxWithTracer, span := otel.Tracer(AuthErrorTracer).Start(ctx.Request.Context(), "Authenticate")
		defer span.End()

		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			err := ce.NewError(span, ce.CodeAuthUnauthenticated, ce.MsgUnauthenticated, errors.New("authorization header is missing"))
			ctx.Error(err)
			ctx.Abort()
			return
		}

		authSlice := strings.Split(authHeader, " ")
		if len(authSlice) != 2 {
			err := ce.NewError(span, ce.CodeAuthTokenMalformed, ce.MsgUnauthenticated, errors.New("authorization header format is invalid"))
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if authType := strings.ToLower(authSlice[0]); authType != "bearer" {
			err := ce.NewError(span, ce.CodeAuthTokenMalformed, ce.MsgUnauthenticated, errors.New("authorization type must be Bearer"))
			ctx.Error(err)
			ctx.Abort()
			return
		}

		token := authSlice[1]
		claim, err := utils.ParseJWTToken(token)
		if err != nil {
			switch {
			case errors.Is(err, ce.ErrTokenExpired):
				err = ce.NewError(span, ce.CodeAuthTokenExpired, ce.MsgUnauthenticated, err)
			case errors.Is(err, ce.ErrTokenMalformed):
				err = ce.NewError(span, ce.CodeAuthTokenMalformed, ce.MsgUnauthenticated, err)
			default:
				err = ce.NewError(span, ce.CodeAuthTokenParsing, ce.MsgInternalServer, err)
			}

			ctx.Error(err)
			ctx.Abort()
			return
		}

		if !utils.IsInTokenAudience(claim.Audience) {
			err := ce.NewError(span, ce.CodeAuthAudienceNotFound, ce.MsgUnauthenticated, errors.New("service not found in token audience"))
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctxWithTracer = context.WithValue(ctxWithTracer, constants.CtxKeyAuthID, claim.AuthID)
		ctxWithTracer = context.WithValue(ctxWithTracer, constants.CtxKeyRoleID, claim.RoleID)
		ctxWithTracer = context.WithValue(ctxWithTracer, constants.CtxKeyIsVerified, claim.IsVerified)

		ctx.Request = ctx.Request.WithContext(ctxWithTracer)
		ctx.Next()
	}
}
