package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/validator"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const authErrorTracer string = "handler.auth"

type AuthHandler struct {
	au        usecases.AuthUsecase
	validator *validator.Validator
	cookie    *services.CookieService
	cfg       *configs.Config
}

func NewAuthHandler(au usecases.AuthUsecase, validator *validator.Validator, cookie *services.CookieService, cfg *configs.Config) *AuthHandler {
	return &AuthHandler{au, validator, cookie, cfg}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Register")
	defer span.End()

	var payload dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to register: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		err := fmt.Errorf("failed to register: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, err))
		return
	}

	data := entities.CreateAuth{
		Email:    payload.Email,
		Password: &payload.Password,
	}
	request := entities.Request{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	authToken, auth, err := h.au.Register(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RegisterResponse{
		Token: authToken.AccessToken,
		Auth:  h.toAuthResponse(*auth),
	}

	h.setCookie(ctx, authToken.SessionToken)
	utils.SetResponse(ctx, "Registered successfully", response, http.StatusCreated)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Login")
	defer span.End()

	var payload dto.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to login: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		err := fmt.Errorf("failed to login: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, err))
		return
	}

	data := entities.GetAuth{
		Email:    payload.Email,
		Password: payload.Password,
	}
	request := entities.Request{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	authToken, auth, err := h.au.Login(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.LoginResponse{
		Token: authToken.AccessToken,
		Auth:  h.toAuthResponse(*auth),
	}

	h.setCookie(ctx, authToken.SessionToken)
	utils.SetResponse(ctx, "Logged in successfully", response, http.StatusOK)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "Logout")
	defer span.End()

	sessionToken, err := ctx.Cookie(constants.CookieKeySessionToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		wErr := fmt.Errorf("failed to logout: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgUnauthenticated, wErr))
		return
	}
	if err != nil {
		wErr := fmt.Errorf("failed to logout: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgInternalServer, wErr))
		return
	}
	if sessionToken == "" {
		wErr := fmt.Errorf("failed to logout: %w", ce.ErrCookieNotFound)
		ctx.Error(ce.NewError(span, ce.CodeSessionNotFound, ce.MsgUnauthenticated, wErr))
		return
	}

	if err := h.au.Logout(ctxWithTracer, sessionToken); err != nil {
		ctx.Error(err)
		return
	}

	h.delCookie(ctx)
	utils.SetResponse(ctx, "", nil, http.StatusNoContent)
}

func (h *AuthHandler) ChangeEmail(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ChangeEmail")
	defer span.End()

	authID, err := utils.CtxGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to change email: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	var payload dto.ChangeEmailRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to change email: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		err := fmt.Errorf("failed to change email: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, err))
		return
	}

	email, err := h.au.ChangeEmail(ctxWithTracer, authID, payload.NewEmail)
	if err != nil {
		ctx.Error(err)
		return
	}

	msg := fmt.Sprintf("Link to confirm email change sent to %s", email)
	utils.SetResponse(ctx, msg, nil, http.StatusOK)
}

func (h *AuthHandler) ConfirmEmailChange(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ConfirmEmailChange")
	defer span.End()

	var params dto.ConfirmEmailChangeRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		wErr := fmt.Errorf("failed to confirm email change: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	token := strings.TrimSpace(params.Token)
	if token == "" {
		err := fmt.Errorf("failed to confirm email change: %w", ce.ErrTokenNotFound)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err))
		return
	}

	sessionToken, err := ctx.Cookie(constants.CookieKeySessionToken)
	if err != nil {
		// non-fatal: trace the failure, but continue
		span.AddEvent(
			"session cookie not found",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
	}

	authToken, auth, err := h.au.ConfirmEmailChange(ctxWithTracer, token, sessionToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	msg := "Email changed successfully"
	if authToken == nil {
		utils.SetResponse(ctx, msg, nil, http.StatusOK)
		return
	}

	response := dto.ConfirmEmailChangeResponse{
		Token: authToken.AccessToken,
		Auth:  h.toAuthResponse(*auth),
	}

	h.setCookie(ctx, authToken.SessionToken)
	utils.SetResponse(ctx, msg, response, http.StatusOK)
}

func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ChangePassword")
	defer span.End()

	authID, err := utils.CtxGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to change password: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	var payload dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to change password: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}

	data := entities.UpdatePassword{
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	}
	if err := h.au.ChangePassword(ctxWithTracer, authID, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Password changed successfully", nil, http.StatusOK)
}

func (h *AuthHandler) ForgotPassword(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ForgotPassword")
	defer span.End()

	var payload dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to forgot password: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		err := fmt.Errorf("failed to forgot password: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, err))
		return
	}

	email, err := h.au.ForgotPassword(ctxWithTracer, payload.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	msg := fmt.Sprintf("Link to reset password sent to %s", email)
	utils.SetResponse(ctx, msg, nil, http.StatusOK)
}

func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ResetPassword")
	defer span.End()

	var payload dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to reset password: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}

	token := strings.TrimSpace(payload.Token)
	if token == "" {
		err := fmt.Errorf("failed to reset password: %w", ce.ErrTokenNotFound)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err))
		return
	}

	data := entities.ResetPassword{
		Token:       token,
		NewPassword: payload.NewPassword,
	}
	if err := h.au.ResetPassword(ctxWithTracer, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Password changed successfully", nil, http.StatusOK)
}

func (h *AuthHandler) ResendVerification(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "ResendVerification")
	defer span.End()

	authID, err := utils.CtxGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to resend verification: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	email, err := h.au.ResendVerification(ctxWithTracer, authID)
	if err != nil {
		ctx.Error(err)
		return
	}

	msg := fmt.Sprintf("Link to verify account sent to %s", email)
	utils.SetResponse(ctx, msg, nil, http.StatusOK)
}

func (h *AuthHandler) VerifyAccount(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "VerifyAccount")
	defer span.End()

	var params dto.VerifyAccountRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		wErr := fmt.Errorf("failed to verify account: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	token := strings.TrimSpace(params.Token)
	if token == "" {
		err := fmt.Errorf("failed to verify account: %w", ce.ErrTokenNotFound)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err))
		return
	}

	sessionToken, err := ctx.Cookie(constants.CookieKeySessionToken)
	if err != nil {
		// non-fatal: trace the failure, but continue
		span.AddEvent(
			"session cookie not found",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
	}

	authToken, auth, err := h.au.VerifyAccount(ctxWithTracer, token, sessionToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	msg := "Account verified successfully"
	if authToken == nil {
		utils.SetResponse(ctx, msg, nil, http.StatusOK)
		return
	}

	response := dto.VerifyAccountResponse{
		Token: authToken.AccessToken,
		Auth:  h.toAuthResponse(*auth),
	}

	h.setCookie(ctx, authToken.SessionToken)
	utils.SetResponse(ctx, msg, response, http.StatusOK)
}

func (h *AuthHandler) RefreshSession(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "RefreshSession")
	defer span.End()

	sessionToken, err := ctx.Cookie(constants.CookieKeySessionToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		wErr := fmt.Errorf("failed to refresh session: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgUnauthenticated, wErr))
		return
	}
	if err != nil {
		wErr := fmt.Errorf("failed to refresh session: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextCookieNotFound, ce.MsgInternalServer, wErr))
		return
	}

	sessionToken = strings.TrimSpace(sessionToken)
	if sessionToken == "" {
		err := fmt.Errorf("failed to refresh session: %w", ce.ErrTokenNotFound)
		ctx.Error(ce.NewError(span, ce.CodeSessionNotFound, ce.MsgUnauthenticated, err))
		return
	}

	authToken, err := h.au.RefreshSession(ctxWithTracer, sessionToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RefreshSessionResponse{
		Token: authToken.AccessToken,
	}

	h.setCookie(ctx, authToken.SessionToken)
	utils.SetResponse(ctx, "Session refreshed successfully", response, http.StatusOK)
}

func (h *AuthHandler) IsEmailRegistered(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "IsEmailRegistered")
	defer span.End()

	var params dto.QueryEmailRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		wErr := fmt.Errorf("failed to query email registration: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	isRegistered, err := h.au.IsEmailRegistered(ctxWithTracer, params.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.QueryEmailResponse{
		IsRegistered: isRegistered,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *AuthHandler) IsResetTokenValid(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "IsResetTokenValid")
	defer span.End()

	var payload dto.QueryTokenRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to query reset token validity: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}

	token := strings.TrimSpace(payload.Token)
	if token == "" {
		err := fmt.Errorf("failed to query reset token validity: %w", ce.ErrTokenNotFound)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err))
		return
	}

	isValid, err := h.au.IsResetTokenValid(ctxWithTracer, token)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.QueryTokenResponse{
		IsValid: isValid,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *AuthHandler) toAuthResponse(auth entities.Auth) dto.AuthResponse {
	return dto.AuthResponse{
		ID:         auth.ID,
		Email:      auth.Email,
		RoleID:     auth.RoleID,
		IsVerified: auth.IsVerified,
		CreatedAt:  auth.CreatedAt,
		UpdatedAt:  auth.UpdatedAt,
	}
}

func (h *AuthHandler) setCookie(ctx *gin.Context, sessionToken string) {
	h.cookie.Set(ctx, constants.CookieKeySessionToken, sessionToken, h.cfg.Auth.TokenDuration.Session, "/", h.cfg.Server.Host)
}

func (h *AuthHandler) delCookie(ctx *gin.Context) {
	h.cookie.Delete(ctx, constants.CookieKeySessionToken, "/", h.cfg.Server.Host)
}
