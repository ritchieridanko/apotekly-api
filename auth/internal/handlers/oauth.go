package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"golang.org/x/oauth2"
)

const OAuthErrorTracer = ce.OAuthHandlerTracer

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
	tracer := OAuthErrorTracer + ": GoogleCallback()"

	var params dto.ReqOAuth
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidParams, ce.ErrMsgInvalidParams, tracer, err)
		ctx.Error(err)
		return
	}

	token, err := h.google.Exchange(ctx, params.Code)
	if err != nil {
		err := ce.NewError(ce.ErrCodeOAuth, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	user, err := utils.GetUserFromGoogle(ctx, token, h.google)
	if err != nil {
		err := ce.NewError(ce.ErrCodeOAuth, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	data := entities.NewOAuth{
		Provider:   constants.OAuthProviderGoogle,
		UID:        user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	authToken, err := h.oau.Authenticate(ctx, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SetSessionCookie(ctx, authToken.SessionToken)
	url := utils.GenerateURLWithTokenQuery("/auth/oauth-callback", authToken.AccessToken)

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oauthHandler) MicrosoftOAuth(ctx *gin.Context) {
	url := h.microsoft.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oauthHandler) MicrosoftCallback(ctx *gin.Context) {
	tracer := OAuthErrorTracer + ": MicrosoftCallback()"

	var params dto.ReqOAuth
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidParams, ce.ErrMsgInvalidParams, tracer, err)
		ctx.Error(err)
		return
	}

	token, err := h.microsoft.Exchange(ctx, params.Code)
	if err != nil {
		err := ce.NewError(ce.ErrCodeOAuth, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	user, err := utils.GetUserFromMicrosoft(ctx, token.AccessToken)
	if err != nil {
		err := ce.NewError(ce.ErrCodeOAuth, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	email := user.Mail
	if email == "" {
		email = user.UserPrincipalName
	}

	data := entities.NewOAuth{
		Provider:   constants.OAuthProviderMicrosoft,
		UID:        user.ID,
		Email:      email,
		IsVerified: true,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	authToken, err := h.oau.Authenticate(ctx, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SetSessionCookie(ctx, authToken.SessionToken)
	url := utils.GenerateURLWithTokenQuery("/auth/oauth-callback", authToken.AccessToken)

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
