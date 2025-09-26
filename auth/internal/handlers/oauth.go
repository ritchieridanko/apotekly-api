package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/oauth2"
)

const oAuthErrorTracer string = "handler.oauth"

type OAuthHandler interface {
	GoogleOAuth(ctx *gin.Context)
	GoogleCallback(ctx *gin.Context)
	MicrosoftOAuth(ctx *gin.Context)
	MicrosoftCallback(ctx *gin.Context)
}

type oauthHandler struct {
	oau       usecases.OAuthUsecase
	au        usecases.AuthUsecase
	google    *oauth2.Config
	microsoft *oauth2.Config
}

func NewOAuthHandler(oau usecases.OAuthUsecase, au usecases.AuthUsecase, google *oauth2.Config, microsoft *oauth2.Config) OAuthHandler {
	return &oauthHandler{oau, au, google, microsoft}
}

func (h *oauthHandler) GoogleOAuth(ctx *gin.Context) {
	url := h.google.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oauthHandler) GoogleCallback(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(oAuthErrorTracer).Start(ctx.Request.Context(), "GoogleCallback")
	defer span.End()

	var params dto.ReqOAuth
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err)
		ctx.Error(err)
		return
	}

	token, err := h.google.Exchange(ctxWithTracer, params.Code)
	if err != nil {
		err := ce.NewError(span, ce.CodeOAuthCodeExchangeFailed, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	user, err := utils.OAuthGoogleGetUserInfo(ctxWithTracer, token, h.google)
	if err != nil {
		err := ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	data := entities.OAuth{
		Provider:   constants.OAuthProviderGoogle,
		UID:        user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	authToken, err := h.oau.Authenticate(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.CookieSetSession(ctx, authToken.SessionToken)
	url := utils.GenerateURLWithTokenQuery("/auth/oauth-callback", authToken.AccessToken)

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oauthHandler) MicrosoftOAuth(ctx *gin.Context) {
	url := h.microsoft.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oauthHandler) MicrosoftCallback(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(oAuthErrorTracer).Start(ctx.Request.Context(), "MicrosoftCallback")
	defer span.End()

	var params dto.ReqOAuth
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err)
		ctx.Error(err)
		return
	}

	token, err := h.microsoft.Exchange(ctxWithTracer, params.Code)
	if err != nil {
		err := ce.NewError(span, ce.CodeOAuthCodeExchangeFailed, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	user, err := utils.OAuthMicrosoftGetUserInfo(ctxWithTracer, token.AccessToken)
	if err != nil {
		err := ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	email := user.Mail
	if email == "" {
		email = user.UserPrincipalName
	}

	data := entities.OAuth{
		Provider:   constants.OAuthProviderMicrosoft,
		UID:        user.ID,
		Email:      email,
		IsVerified: true,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	authToken, err := h.oau.Authenticate(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.CookieSetSession(ctx, authToken.SessionToken)
	url := utils.GenerateURLWithTokenQuery("/auth/oauth-callback", authToken.AccessToken)

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
