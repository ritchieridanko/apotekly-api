package usecases

import (
	"context"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/caches"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/email"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

// TODO (1): Send email verification

const AuthErrorTracer = ce.AuthUsecaseTracer

type AuthUsecase interface {
	Register(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
	Login(ctx context.Context, data *entities.GetAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
	Logout(ctx context.Context, sessionToken string) (err error)
	ChangeEmail(ctx context.Context, authID int64, newEmail string) (err error)
	ChangePassword(ctx context.Context, authID int64, data *entities.NewPassword) (err error)
	ForgotPassword(ctx context.Context, email string) (err error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err error)
	IsPasswordResetTokenValid(ctx context.Context, token string) (exists bool, err error)
	RefreshSession(ctx context.Context, sessionToken string) (token *entities.AuthToken, err error)
}

type authUsecase struct {
	ar    repos.AuthRepo
	tx    dbtx.TxManager
	su    SessionUsecase
	oau   OAuthUsecase
	cache caches.Cache
	email email.EmailService
}

func NewAuthUsecase(ar repos.AuthRepo, tx dbtx.TxManager, su SessionUsecase, oau OAuthUsecase, cache caches.Cache, email email.EmailService) AuthUsecase {
	return &authUsecase{ar, tx, su, oau, cache, email}
}

func (u *authUsecase) Register(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	tracer := AuthErrorTracer + ": Register()"

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.GetSessionDuration()) * time.Minute

	var token entities.AuthToken
	err := u.tx.ReturnError(ctx, func(ctx context.Context) error {
		normalizedEmail := utils.NormalizeString(data.Email)

		exists, err := u.ar.IsEmailRegistered(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(ce.ErrCodeConflict, ce.ErrMsgEmailRegistered, tracer, ce.ErrEmailAlreadyRegistered)
		}

		hashedPassword, err := utils.HashPassword(data.Password)
		if err != nil {
			return ce.NewError(ce.ErrCodePassHashing, ce.ErrMsgInternalServer, tracer, err)
		}

		newData := entities.NewAuth{
			Email:    normalizedEmail,
			Password: hashedPassword,
			Role:     constants.RoleCustomer,
		}

		authID, err := u.ar.Create(ctx, &newData)
		if err != nil {
			return err
		}

		sessionToken := utils.GenerateRandomToken()
		accessToken, err := utils.GenerateJWTToken(authID, newData.Role, false)
		if err != nil {
			return ce.NewError(ce.ErrCodeToken, ce.ErrMsgInternalServer, tracer, err)
		}

		sessionData := entities.NewSession{
			AuthID:    authID,
			Token:     sessionToken,
			UserAgent: request.UserAgent,
			IPAddress: request.IPAddress,
			ExpiresAt: now.Add(sessionDuration),
		}

		if err := u.su.NewSessionOnRegister(ctx, &sessionData); err != nil {
			return err
		}

		token = entities.AuthToken{
			AccessToken:  accessToken,
			SessionToken: sessionToken,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// TODO (1): Send email verification

	return &token, nil
}

func (u *authUsecase) Login(ctx context.Context, data *entities.GetAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	tracer := AuthErrorTracer + ": Login()"

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.GetSessionDuration()) * time.Minute
	lockDuration := time.Duration(config.GetAuthLockDuration()) * time.Minute

	normalizedEmail := utils.NormalizeString(data.Email)
	auth, err := u.ar.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}
	if auth.Password == nil {
		return nil, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, ce.ErrOAuthRegularLogin)
	}
	if auth.LockedUntil != nil && auth.LockedUntil.After(now) {
		return nil, ce.NewError(ce.ErrCodeLocked, ce.ErrMsgAccountLocked, tracer, ce.ErrAccountLocked)
	}

	totalFailedAuthKey := utils.GenerateDynamicRedisKey(constants.RedisKeyTotalFailedAuth, auth.ID)
	if err := utils.ValidatePassword(*auth.Password, data.Password); err != nil {
		shouldBeLocked, cacheErr := u.cache.ShouldAccountBeLocked(ctx, totalFailedAuthKey)
		if cacheErr != nil {
			return nil, cacheErr
		}
		if shouldBeLocked {
			if err := u.ar.LockAccount(ctx, auth.ID, now.Add(lockDuration)); err != nil {
				return nil, err
			}
			if err := u.cache.Del(ctx, totalFailedAuthKey); err != nil {
				return nil, err
			}
			return nil, ce.NewError(ce.ErrCodeLocked, ce.ErrMsgAccountLocked, tracer, ce.ErrAccountLocked)
		}

		return nil, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, err)
	}

	if err := u.cache.Del(ctx, totalFailedAuthKey); err != nil {
		return nil, err
	}

	sessionToken := utils.GenerateRandomToken()
	accessToken, err := utils.GenerateJWTToken(auth.ID, auth.Role, auth.IsVerified)
	if err != nil {
		return nil, ce.NewError(ce.ErrCodeToken, ce.ErrMsgInternalServer, tracer, err)
	}

	sessionData := entities.NewSession{
		AuthID:    auth.ID,
		Token:     sessionToken,
		UserAgent: request.UserAgent,
		IPAddress: request.IPAddress,
		ExpiresAt: now.Add(sessionDuration),
	}

	_, err = u.su.NewSession(ctx, &sessionData)
	if err != nil {
		return nil, err
	}

	token := entities.AuthToken{
		AccessToken:  accessToken,
		SessionToken: sessionToken,
	}

	return &token, nil
}

func (u *authUsecase) Logout(ctx context.Context, sessionToken string) error {
	return u.su.RevokeSession(ctx, sessionToken)
}

func (u *authUsecase) ChangeEmail(ctx context.Context, authID int64, newEmail string) error {
	tracer := AuthErrorTracer + ": ChangeEmail()"

	return u.tx.ReturnError(ctx, func(ctx context.Context) error {
		// Email cannot be changed if registered with oauth
		exists, err := u.oau.IsAuthRegistered(ctx, authID)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgOAuthEmailChange, tracer, ce.ErrOAuthEmailChange)
		}

		normalizedEmail := utils.NormalizeString(newEmail)
		exists, err = u.ar.IsEmailRegistered(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(ce.ErrCodeConflict, ce.ErrMsgEmailRegistered, tracer, ce.ErrEmailAlreadyRegistered)
		}

		return u.ar.UpdateEmail(ctx, authID, normalizedEmail)
	})
}

func (u *authUsecase) ChangePassword(ctx context.Context, authID int64, data *entities.NewPassword) error {
	tracer := AuthErrorTracer + ": ChangePassword()"

	return u.tx.ReturnError(ctx, func(ctx context.Context) error {
		auth, err := u.ar.GetByID(ctx, authID)
		if err != nil {
			return err
		}
		if auth.Password == nil {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgOAuthPasswordChange, tracer, ce.ErrOAuthPasswordChange)
		}
		if err := utils.ValidatePassword(*auth.Password, data.OldPassword); err != nil {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidPassword, tracer, err)
		}

		hashedNewPassword, err := utils.HashPassword(data.NewPassword)
		if err != nil {
			return ce.NewError(ce.ErrCodePassHashing, ce.ErrMsgInternalServer, tracer, err)
		}

		return u.ar.UpdatePassword(ctx, auth.ID, hashedNewPassword)
	})
}

func (u *authUsecase) ForgotPassword(ctx context.Context, email string) error {
	normalizedEmail := utils.NormalizeString(email)
	auth, err := u.ar.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		return err
	}
	if auth.Password == nil {
		return nil
	}

	token := utils.GenerateRandomToken()
	tokenDuration := time.Duration(config.GetPasswordResetTokenDuration()) * time.Minute
	if err := u.cache.NewOrReplacePasswordResetToken(ctx, auth.ID, token, tokenDuration); err != nil {
		return err
	}

	return u.email.SendPasswordResetToken(auth.Email, token)
}

func (u *authUsecase) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	normalizedEmail := utils.NormalizeString(email)
	return u.ar.IsEmailRegistered(ctx, normalizedEmail)
}

func (u *authUsecase) IsPasswordResetTokenValid(ctx context.Context, token string) (bool, error) {
	tokenKey := utils.GenerateDynamicRedisKey(constants.RedisKeyPasswordResetToken, token)
	return u.cache.Has(ctx, tokenKey)
}

func (u *authUsecase) RefreshSession(ctx context.Context, sessionToken string) (*entities.AuthToken, error) {
	tracer := AuthErrorTracer + ": RefreshSession()"

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.GetSessionDuration()) * time.Minute

	var token entities.AuthToken
	err := u.tx.ReturnError(ctx, func(ctx context.Context) error {
		session, err := u.su.GetSession(ctx, sessionToken)
		if err != nil {
			return err
		}
		if !session.ExpiresAt.After(now) {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, ce.ErrSessionExpired)
		}
		if session.RevokedAt != nil {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, ce.ErrSessionRevoked)
		}

		auth, err := u.ar.GetByID(ctx, session.AuthID)
		if err != nil {
			return err
		}

		newSessionToken := utils.GenerateRandomToken()
		newAccessToken, err := utils.GenerateJWTToken(auth.ID, auth.Role, auth.IsVerified)
		if err != nil {
			return ce.NewError(ce.ErrCodeToken, ce.ErrMsgInternalServer, tracer, err)
		}

		newSessionData := entities.ReissueSession{
			AuthID:    auth.ID,
			ParentID:  session.ID,
			Token:     newSessionToken,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			ExpiresAt: now.Add(sessionDuration),
		}

		_, err = u.su.RenewSession(ctx, &newSessionData)
		if err != nil {
			return err
		}

		token = entities.AuthToken{
			AccessToken:  newAccessToken,
			SessionToken: newSessionToken,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &token, nil
}
