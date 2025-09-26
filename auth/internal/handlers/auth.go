package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
)

const authErrorTracer string = "handler.auth"

type AuthHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	ChangeEmail(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ResendVerification(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	IsEmailRegistered(ctx *gin.Context)
	IsResetTokenValid(ctx *gin.Context)
	RefreshSession(ctx *gin.Context)
}

type authHandler struct {
	au usecases.AuthUsecase
}

func NewAuthHandler(au usecases.AuthUsecase) AuthHandler {
	return &authHandler{au}
}

func (h *authHandler) Register(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Register")
	defer span.End()

	var payload dto.ReqRegister
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
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

	token, err := h.au.Register(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: token.AccessToken,
	}

	utils.CookieSetSession(ctx, token.SessionToken)
	utils.SetResponse(ctx, "Registered successfully.", response, http.StatusCreated)
}

func (h *authHandler) Login(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Login")
	defer span.End()

	var payload dto.ReqLogin
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
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

	token, err := h.au.Login(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: token.AccessToken,
	}

	utils.CookieSetSession(ctx, token.SessionToken)
	utils.SetResponse(ctx, "Logged in successfully. Welcome back!", response, http.StatusOK)
}

func (h *authHandler) Logout(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Logout")
	defer span.End()

	token, err := ctx.Cookie(constants.CookieKeySessionToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		err := ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgUnauthenticated, err)
		ctx.Error(err)
		return
	}
	if err != nil {
		err := ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}
	if token == "" {
		err := ce.NewError(span, ce.CodeSessionNotFound, ce.MsgUnauthenticated, errors.New("session token not found in cookie"))
		ctx.Error(err)
		return
	}

	if err := h.au.Logout(ctxWithTracer, token); err != nil {
		ctx.Error(err)
		return
	}

	utils.CookieDelSession(ctx)
	utils.SetResponse(ctx, "", nil, http.StatusNoContent)
}

func (h *authHandler) ChangeEmail(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ChangeEmail")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqChangeEmail
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if err := h.au.ChangeEmail(ctxWithTracer, authID, payload.NewEmail); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Email changed successfully.", nil, http.StatusOK)
}

func (h *authHandler) ChangePassword(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ChangePassword")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqChangePassword
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	data := entities.PasswordChange{
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}

	if err := h.au.ChangePassword(ctxWithTracer, authID, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Password changed successfully.", nil, http.StatusOK)
}

func (h *authHandler) ForgotPassword(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ForgotPassword")
	defer span.End()

	var payload dto.ReqForgotPassword
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if err := h.au.ForgotPassword(ctxWithTracer, payload.Email); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Link sent. Please check your email!", nil, http.StatusOK)
}

func (h *authHandler) ResetPassword(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ResetPassword")
	defer span.End()

	var payload dto.ReqResetPassword
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	token := strings.TrimSpace(payload.Token)
	if token == "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, errors.New("reset token not found in payload"))
		ctx.Error(err)
		return
	}

	data := entities.NewPassword{
		Token:       token,
		NewPassword: payload.NewPassword,
	}

	if err := h.au.ResetPassword(ctxWithTracer, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Password changed successfully.", nil, http.StatusOK)
}

func (h *authHandler) ResendVerification(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ResendVerification")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	if err := h.au.ResendVerification(ctxWithTracer, authID); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Link sent. Please check your email!", nil, http.StatusOK)
}

func (h *authHandler) VerifyEmail(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "VerifyEmail")
	defer span.End()

	var params dto.ReqVerifyEmail
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err)
		ctx.Error(err)
		return
	}

	token := strings.TrimSpace(params.Token)
	if token == "" {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, errors.New("verification token not found in params"))
		ctx.Error(err)
		return
	}

	if err := h.au.VerifyEmail(ctxWithTracer, token); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Email verified successfully.", nil, http.StatusOK)
}

func (h *authHandler) IsEmailRegistered(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "IsEmailRegistered")
	defer span.End()

	var params dto.ReqQueryEmail
	if err := ctx.ShouldBindQuery(&params); err != nil {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err)
		ctx.Error(err)
		return
	}

	isRegistered, err := h.au.IsEmailRegistered(ctxWithTracer, params.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespQueryEmail{
		IsRegistered: isRegistered,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *authHandler) IsResetTokenValid(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "IsResetTokenValid")
	defer span.End()

	var payload dto.ReqQueryToken
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	token := strings.TrimSpace(payload.Token)
	if token == "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, errors.New("reset token not found in payload"))
		ctx.Error(err)
		return
	}

	isValid, err := h.au.IsResetTokenValid(ctxWithTracer, token)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespQueryToken{
		IsValid: isValid,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *authHandler) RefreshSession(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "RefreshSession")
	defer span.End()

	token, err := ctx.Cookie(constants.CookieKeySessionToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		err := ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgUnauthenticated, err)
		ctx.Error(err)
		return
	}
	if err != nil {
		err := ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}
	if token == "" {
		err := ce.NewError(span, ce.CodeSessionNotFound, ce.MsgUnauthenticated, errors.New("session token not found in cookie"))
		ctx.Error(err)
		return
	}

	authToken, err := h.au.RefreshSession(ctxWithTracer, token)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: authToken.AccessToken,
	}

	utils.CookieSetSession(ctx, authToken.SessionToken)
	utils.SetResponse(ctx, "Session refreshed successfully.", response, http.StatusOK)
}
