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

const authErrorTracer string = "middleware.auth"

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Authenticate")
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

		claim, err := utils.ParseJWTToken(authSlice[1])
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

func RequireVerified() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "RequireVerified")
		defer span.End()

		value := ctxWithTracer.Value(constants.CtxKeyIsVerified)
		isVerified, ok := value.(bool)
		if !ok || !isVerified {
			err := ce.NewError(span, ce.CodeAuthNotVerified, "Please verify your email first!", errors.New("account not verified"))
			ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
