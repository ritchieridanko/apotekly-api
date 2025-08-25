package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const AuthErrorTracer = ce.AuthHandlerTracer

type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshSession(ctx *gin.Context)
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
		newToken, err = h.au.RefreshSession(ctx, token)
		if err != nil {
			newToken, err = h.au.Login(ctx, &data, &request)
		}
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

func (h *authHandler) Logout(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": Logout()"

	token, err := ctx.Cookie(constants.CookieKeySessionToken)
	if errors.Is(err, http.ErrNoCookie) {
		err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, err)
		ctx.Error(err)
		return
	}
	if err != nil {
		err := ce.NewError(ce.ErrCodeContext, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}
	if token == "" {
		err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, ce.ErrTokenEmpty)
		ctx.Error(err)
		return
	}

	if err := h.au.Logout(ctx, token); err != nil {
		ctx.Error(err)
		return
	}

	utils.DeleteSessionCookie(ctx)
	utils.SetResponse(ctx, "logged out successfully", nil, http.StatusNoContent)
}

func (h *authHandler) RefreshSession(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": RefreshSession()"

	token, err := ctx.Cookie(constants.CookieKeySessionToken)
	if errors.Is(err, http.ErrNoCookie) {
		err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, err)
		ctx.Error(err)
		return
	}
	if err != nil {
		err := ce.NewError(ce.ErrCodeContext, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}
	if token == "" {
		err := ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUnauthenticated, tracer, ce.ErrTokenEmpty)
		ctx.Error(err)
		return
	}

	newToken, err := h.au.RefreshSession(ctx, token)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: newToken.AccessToken,
	}

	utils.SetSessionCookie(ctx, newToken.SessionToken)
	utils.SetResponse(ctx, "session refreshed successfully", response, http.StatusOK)
}
