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
}

type oauthHandler struct {
	oau    usecases.OAuthUsecase
	au     usecases.AuthUsecase
	google *oauth2.Config
}

func NewOAuthHandler(oau usecases.OAuthUsecase, au usecases.AuthUsecase, google *oauth2.Config) OAuthHandler {
	return &oauthHandler{oau, au, google}
}

func (h *oauthHandler) GoogleOAuth(ctx *gin.Context) {
	url := h.google.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oauthHandler) GoogleCallback(ctx *gin.Context) {
	tracer := OAuthErrorTracer + ": GoogleCallback()"

	var params dto.ReqOAuthByGoogle
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

	userInfo, err := utils.GetUserInfoFromGoogle(ctx, token, h.google)
	if err != nil {
		err := ce.NewError(ce.ErrCodeOAuth, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	data := entities.NewOAuth{
		Provider:   constants.OAuthProviderGoogle,
		UID:        userInfo.ID,
		Email:      userInfo.Email,
		IsVerified: userInfo.IsVerified,
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
	url := utils.GenerateURLWithTokenQuery("/oauth/google/callback", authToken.AccessToken)

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
