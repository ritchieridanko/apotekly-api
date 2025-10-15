package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"go.opentelemetry.io/otel"
	"golang.org/x/oauth2"
)

const oAuthErrorTracer string = "handler.oauth"

type OAuthHandler struct {
	oau       usecases.OAuthUsecase
	au        usecases.AuthUsecase
	google    *oauth2.Config
	microsoft *oauth2.Config
	cookie    *services.CookieService
	cfg       *configs.Config
}

func NewOAuthHandler(
	oau usecases.OAuthUsecase,
	au usecases.AuthUsecase,
	google, microsoft *oauth2.Config,
	cookie *services.CookieService,
	cfg *configs.Config,
) *OAuthHandler {
	return &OAuthHandler{oau, au, google, microsoft, cookie, cfg}
}

func (h *OAuthHandler) GoogleOAuth(ctx *gin.Context) {
	url := h.google.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *OAuthHandler) GoogleCallback(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(oAuthErrorTracer).Start(ctx.Request.Context(), "GoogleCallback")
	defer span.End()

	var params dto.AuthenticateRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		wErr := fmt.Errorf("failed to handle google callback: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	token, err := h.google.Exchange(ctxWithTracer, params.Code)
	if err != nil {
		wErr := fmt.Errorf("failed to handle google callback: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeOAuthCodeExchangeFailed, ce.MsgInternalServer, wErr))
		return
	}

	user, err := h.googleGetUserInfo(ctxWithTracer, token, h.google)
	if err != nil {
		ctx.Error(err)
		return
	}

	data := entities.OAuth{
		Provider:   constants.OAuthProviderGoogle,
		UID:        user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}
	request := entities.Request{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	sessionToken, exchangeCode, err := h.oau.Authenticate(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	h.setCookie(ctx, sessionToken)

	url := h.setRedirectURL(exchangeCode)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *OAuthHandler) MicrosoftOAuth(ctx *gin.Context) {
	url := h.microsoft.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *OAuthHandler) MicrosoftCallback(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(oAuthErrorTracer).Start(ctx.Request.Context(), "MicrosoftCallback")
	defer span.End()

	var params dto.AuthenticateRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		wErr := fmt.Errorf("failed to handle microsoft callback: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	token, err := h.microsoft.Exchange(ctxWithTracer, params.Code)
	if err != nil {
		wErr := fmt.Errorf("failed to handle microsoft callback: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeOAuthCodeExchangeFailed, ce.MsgInternalServer, wErr))
		return
	}

	user, err := h.microsoftGetUserInfo(ctxWithTracer, token.AccessToken)
	if err != nil {
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
	request := entities.Request{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	sessionToken, exchangeCode, err := h.oau.Authenticate(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	h.setCookie(ctx, sessionToken)

	url := h.setRedirectURL(exchangeCode)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *OAuthHandler) OAuthExchange(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(oAuthErrorTracer).Start(ctx.Request.Context(), "OAuthExchange")
	defer span.End()

	var params dto.ExchangeCodeRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		wErr := fmt.Errorf("failed to handle oauth code exchange: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	code := strings.TrimSpace(params.Code)
	if code == "" {
		err := fmt.Errorf("failed to handle oauth code exchange: %w", ce.ErrOAuthCodeNotFound)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err))
		return
	}

	auth, accessToken, err := h.oau.ExchangeCode(ctxWithTracer, code)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.ExchangeCodeResponse{
		Token: accessToken,
		Auth:  h.toAuthResponse(*auth),
	}

	utils.SetResponse(ctx, "Code exchanged successfully", response, http.StatusOK)
}

func (h *OAuthHandler) googleGetUserInfo(ctx context.Context, token *oauth2.Token, cfg *oauth2.Config) (*dto.GoogleUser, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "googleGetUserInfo")
	defer span.End()

	client := cfg.Client(ctx, token)

	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		wErr := fmt.Errorf("failed to fetch google user info: %w", err)
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, wErr)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to fetch google user info: %w", fmt.Errorf("response status '%s'", response.Status))
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, err)
	}

	var user dto.GoogleUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		wErr := fmt.Errorf("failed to fetch google user info: %w", err)
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, wErr)
	}

	return &user, nil
}

func (h *OAuthHandler) microsoftGetUserInfo(ctx context.Context, accessToken string) (*dto.MicrosoftUser, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "microsoftGetUserInfo")
	defer span.End()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch microsoft user info: %w", err)
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, wErr)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch microsoft user info: %w", err)
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, wErr)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to fetch microsoft user info: %w", fmt.Errorf("response status '%s'", response.Status))
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, err)
	}

	var user dto.MicrosoftUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		wErr := fmt.Errorf("failed to fetch microsoft user info: %w", err)
		return nil, ce.NewError(span, ce.CodeOAuthGetUserInfoFailed, ce.MsgInternalServer, wErr)
	}

	return &user, nil
}

func (h *OAuthHandler) toAuthResponse(auth entities.Auth) dto.AuthResponse {
	return dto.AuthResponse{
		ID:         auth.ID,
		Email:      auth.Email,
		RoleID:     auth.RoleID,
		IsVerified: auth.IsVerified,
		CreatedAt:  auth.CreatedAt,
		UpdatedAt:  auth.UpdatedAt,
	}
}

func (h *OAuthHandler) setCookie(ctx *gin.Context, sessionToken string) {
	h.cookie.Set(ctx, constants.CookieKeySessionToken, sessionToken, h.cfg.Auth.TokenDuration.Session, "/", h.cfg.Server.Host)
}

func (h *OAuthHandler) setRedirectURL(code string) string {
	return fmt.Sprintf("%s/auth/oauth-callback?code=%s", h.cfg.Client.BaseURL, code)
}
