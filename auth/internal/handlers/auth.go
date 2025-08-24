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
)

// TODO (1): Refresh session

const AuthErrorTracer = ce.AuthHandlerTracer

type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authHandler struct {
	au usecases.AuthUsecase
}

func NewAuthHandler(au usecases.AuthUsecase) AuthHandler {
	return &authHandler{au}
}

func (h *authHandler) Register(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": Register()"

	var payload dto.ReqRegister
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	data := entities.NewAuth{
		Email:    payload.Email,
		Password: payload.Password,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	token, err := h.au.Register(ctx, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: token.AccessToken,
	}

	utils.SetSessionCookie(ctx, token.SessionToken)
	utils.SetResponse(ctx, "registered successfully", response, http.StatusCreated)
}

func (h *authHandler) Login(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": Login()"

	var payload dto.ReqLogin
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	data := entities.GetAuth{
		Email:    payload.Email,
		Password: payload.Password,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	var newToken *entities.AuthToken
	token, err := ctx.Cookie(constants.CookieKeySessionToken)
	if err == nil && token != "" {
		// TODO (1): Refresh session
	} else {
		newToken, err = h.au.Login(ctx, &data, &request)
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: newToken.AccessToken,
	}

	utils.SetSessionCookie(ctx, newToken.SessionToken)
	utils.SetResponse(ctx, "logged in successfully", response, http.StatusOK)
}
