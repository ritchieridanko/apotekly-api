package handlers

import (
	"errors"
	"net/http"
	"strings"

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
	ChangeEmail(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	IsEmailRegistered(ctx *gin.Context)
	IsPasswordResetTokenValid(ctx *gin.Context)
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

	token, err := h.au.Login(ctx, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: token.AccessToken,
	}

	utils.SetSessionCookie(ctx, token.SessionToken)
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
	utils.SetResponse(ctx, "", nil, http.StatusNoContent)
}

func (h *authHandler) ChangeEmail(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": ChangeEmail()"

	authID, err := utils.GetAuthIDFromContext(ctx)
	if err != nil {
		err := ce.NewError(ce.ErrCodeContext, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqEmailChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	if err := h.au.ChangeEmail(ctx, authID, payload.NewEmail); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "email changed successfully", nil, http.StatusOK)
}

func (h *authHandler) ChangePassword(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": ChangePassword()"

	authID, err := utils.GetAuthIDFromContext(ctx)
	if err != nil {
		err := ce.NewError(ce.ErrCodeContext, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqPasswordChange
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	data := entities.NewPassword{
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}

	if err := h.au.ChangePassword(ctx, authID, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "password changed successfully", nil, http.StatusOK)
}

func (h *authHandler) ForgotPassword(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": ForgotPassword()"

	var payload dto.ReqForgotPassword
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	if err := h.au.ForgotPassword(ctx, payload.Email); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "please check your email", nil, http.StatusOK)
}

func (h *authHandler) ResetPassword(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": ResetPassword()"

	var payload dto.ReqPasswordReset
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	token := strings.TrimSpace(payload.Token)
	if token == "" {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, ce.ErrTokenEmpty)
		ctx.Error(err)
		return
	}

	data := entities.PasswordReset{
		Token:       token,
		NewPassword: payload.NewPassword,
	}

	if err := h.au.ResetPassword(ctx, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "password changed successfully", nil, http.StatusOK)
}

func (h *authHandler) IsEmailRegistered(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": IsEmailRegistered()"

	var params dto.ReqEmailCheckQuery
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidParams, ce.ErrMsgInvalidParams, tracer, err)
		ctx.Error(err)
		return
	}

	isRegistered, err := h.au.IsEmailRegistered(ctx, params.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespEmailCheckQuery{
		IsRegistered: isRegistered,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *authHandler) IsPasswordResetTokenValid(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": IsPasswordResetTokenValid()"

	var payload dto.ReqTokenCheckQuery
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	token := strings.TrimSpace(payload.Token)
	if token == "" {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, ce.ErrTokenEmpty)
		ctx.Error(err)
		return
	}

	isValid, err := h.au.IsPasswordResetTokenValid(ctx, token)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespTokenCheckQuery{
		IsValid: isValid,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
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
