package middlewares

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"go.opentelemetry.io/otel"
)

const authErrorTracer string = "middleware.auth"

type AuthMiddleware struct {
	jwt     *services.JWTService
	appName string
}

func NewAuthMiddleware(jwt *services.JWTService, appName string) *AuthMiddleware {
	return &AuthMiddleware{jwt, appName}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Authenticate")
		defer span.End()

		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if len(authHeader) == 0 {
			wErr := fmt.Errorf("failed to authenticate: %w", errors.New("authorization header is missing"))
			ctx.Error(ce.NewError(span, ce.CodeAuthUnauthenticated, ce.MsgUnauthenticated, wErr))
			ctx.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			wErr := fmt.Errorf("failed to authenticate: %w", errors.New("invalid authorization format"))
			ctx.Error(ce.NewError(span, ce.CodeAuthTokenMalformed, ce.MsgUnauthenticated, wErr))
			ctx.Abort()
			return
		}

		claim, err := m.jwt.Parse(authParts[1])
		if err != nil {
			wErr := fmt.Errorf("failed to authenticate: %w", err)
			switch {
			case errors.Is(err, ce.ErrTokenExpired):
				err = ce.NewError(span, ce.CodeAuthTokenExpired, ce.MsgUnauthenticated, wErr)
			case errors.Is(err, ce.ErrTokenMalformed):
				err = ce.NewError(span, ce.CodeAuthTokenMalformed, ce.MsgUnauthenticated, wErr)
			case errors.Is(err, ce.ErrInvalidTokenClaim):
				err = ce.NewError(span, ce.CodeInvalidTokenClaim, ce.MsgUnauthenticated, wErr)
			default:
				err = ce.NewError(span, ce.CodeAuthTokenParsing, ce.MsgInternalServer, wErr)
			}

			ctx.Error(err)
			ctx.Abort()
			return
		}

		if !m.jwt.ValidateAudience(m.appName) {
			wErr := fmt.Errorf("failed to authenticate: %w", errors.New("service not supported"))
			ctx.Error(ce.NewError(span, ce.CodeAuthAudienceNotFound, ce.MsgUnauthenticated, wErr))
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

func (m *AuthMiddleware) Authorize(roleID int16) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Authorize")
		defer span.End()

		value := ctxWithTracer.Value(constants.CtxKeyRoleID)
		r, ok := value.(int16)
		if !ok || r != roleID {
			wErr := fmt.Errorf("failed to authorize: %w", errors.New("role unauthorized"))
			ctx.Error(ce.NewError(span, ce.CodeRoleUnauthorized, "Unauthorized", wErr))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *AuthMiddleware) RequireVerified() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "RequireVerified")
		defer span.End()

		value := ctxWithTracer.Value(constants.CtxKeyIsVerified)
		isVerified, ok := value.(bool)
		if !ok || !isVerified {
			wErr := fmt.Errorf("failed to require verification: %w", errors.New("account not verified"))
			ctx.Error(ce.NewError(span, ce.CodeAuthNotVerified, "Please verify your email first!", wErr))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
