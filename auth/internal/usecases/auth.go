package usecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/email"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
)

const authErrorTracer string = "usecase.auth"

type AuthUsecase interface {
	Register(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
	Login(ctx context.Context, data *entities.GetAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
	Logout(ctx context.Context, sessionToken string) (err error)
	ChangeEmail(ctx context.Context, authID int64, newEmail string) (err error)
	ChangePassword(ctx context.Context, authID int64, data *entities.PasswordChange) (err error)
	ForgotPassword(ctx context.Context, email string) (err error)
	ResetPassword(ctx context.Context, data *entities.NewPassword) (err error)
	ResendVerification(ctx context.Context, authID int64) (err error)
	VerifyEmail(ctx context.Context, token string) (err error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err error)
	IsResetTokenValid(ctx context.Context, token string) (exists bool, err error)
	RefreshSession(ctx context.Context, sessionToken string) (token *entities.AuthToken, err error)
}

type authUsecase struct {
	ar    repos.AuthRepo
	su    SessionUsecase
	tx    db.TxManager
	cache cache.CacheService
	email email.EmailService
}

func NewAuthUsecase(
	ar repos.AuthRepo,
	su SessionUsecase,
	tx db.TxManager,
	cache cache.CacheService,
	email email.EmailService,
) AuthUsecase {
	return &authUsecase{ar, su, tx, cache, email}
}

func (u *authUsecase) Register(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Register")
	defer span.End()

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.AuthGetSessionDuration()) * time.Minute

	var normalizedEmail string
	var authID int64
	var token entities.AuthToken
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {
		normalizedEmail = utils.NormalizeString(data.Email)
		exists, err := u.ar.IsEmailRegistered(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(span, ce.CodeAuthEmailConflict, "Email is already registered.", errors.New("registration email conflict"))
		}

		hashedPassword, err := utils.HashPassword(data.Password)
		if err != nil {
			return ce.NewError(span, ce.CodePasswordHashingFailed, ce.MsgInternalServer, err)
		}

		newData := entities.NewAuth{
			Email:    normalizedEmail,
			Password: hashedPassword,
			Role:     constants.RoleCustomer,
		}

		authID, err = u.ar.Create(ctx, &newData)
		if err != nil {
			return err
		}

		sessionToken := utils.GenerateRandomToken()
		accessToken, err := utils.GenerateJWTToken(authID, newData.Role, false)
		if err != nil {
			return ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, err)
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

	verificationToken := utils.GenerateRandomToken()
	tokenDuration := time.Duration(config.AuthGetVerifyTokenDuration()) * time.Minute
	if err := u.cache.NewOrReplaceVerificationToken(ctx, authID, verificationToken, tokenDuration); err != nil {
		log.Println("WARNING -> failed to store verification token in cache after registration:", err.Error())
		return &token, nil
	}

	if err := u.email.SendWelcomeMessage(ctx, normalizedEmail, verificationToken); err != nil {
		log.Println("WARNING -> failed to send welcome email after registration:", err.Error())
	}

	return &token, nil
}

func (u *authUsecase) Login(ctx context.Context, data *entities.GetAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Login")
	defer span.End()

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.AuthGetSessionDuration()) * time.Minute
	lockDuration := time.Duration(config.AuthGetLockDuration()) * time.Minute

	normalizedEmail := utils.NormalizeString(data.Email)
	auth, err := u.ar.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}
	if auth.Password == nil {
		return nil, ce.NewError(span, ce.CodeOAuthRegularLogin, ce.MsgInvalidCredentials, errors.New("email registered as oauth"))
	}
	if auth.LockedUntil != nil && auth.LockedUntil.After(now) {
		return nil, ce.NewError(span, ce.CodeAuthLocked, "Your account is locked. Please try again later!", errors.New("account locked"))
	}

	totalFailedAuthKey := utils.CacheCreateKey(constants.CacheKeyTotalFailedAuth, auth.ID)
	if err := utils.ValidatePassword(*auth.Password, data.Password); err != nil {
		shouldBeLocked, cacheErr := u.cache.ShouldLockAccount(ctx, auth.ID)
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
			return nil, ce.NewError(span, ce.CodeAuthLocked, "Account locked due to multiple failed attempts.", errors.New("multiple failed authentication attempts"))
		}

		return nil, ce.NewError(span, ce.CodeAuthWrongPassword, ce.MsgInvalidCredentials, err)
	}
	if err := u.cache.Del(ctx, totalFailedAuthKey); err != nil {
		return nil, err
	}

	sessionToken := utils.GenerateRandomToken()
	accessToken, err := utils.GenerateJWTToken(auth.ID, auth.Role, auth.IsVerified)
	if err != nil {
		return nil, ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, err)
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
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Logout")
	defer span.End()

	return u.su.RevokeSession(ctx, sessionToken)
}

func (u *authUsecase) ChangeEmail(ctx context.Context, authID int64, newEmail string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ChangeEmail")
	defer span.End()

	return u.tx.WithTx(ctx, func(ctx context.Context) error {
		auth, err := u.ar.GetByID(ctx, authID)
		if err != nil {
			return err
		}
		if auth.Password == nil {
			// email cannot be changed if registered with oauth
			return ce.NewError(span, ce.CodeOAuthEmailChange, "OAuth account cannot change email.", errors.New("email change with oauth account"))
		}

		normalizedEmail := utils.NormalizeString(newEmail)
		exists, err := u.ar.IsEmailRegistered(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(span, ce.CodeAuthEmailConflict, "Email is already registered.", errors.New("registration email conflict"))
		}

		return u.ar.UpdateEmail(ctx, auth.ID, normalizedEmail)
	})
}

func (u *authUsecase) ChangePassword(ctx context.Context, authID int64, data *entities.PasswordChange) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ChangePassword")
	defer span.End()

	return u.tx.WithTx(ctx, func(ctx context.Context) error {
		auth, err := u.ar.GetByID(ctx, authID)
		if err != nil {
			return err
		}
		if auth.Password == nil {
			// password cannot be changed if registered with oauth
			return ce.NewError(span, ce.CodeOAuthPasswordChange, "OAuth account cannot change password.", errors.New("password change with oauth account"))
		}
		if err := utils.ValidatePassword(*auth.Password, data.OldPassword); err != nil {
			return ce.NewError(span, ce.CodeAuthWrongPassword, "Invalid old password.", err)
		}

		hashedNewPassword, err := utils.HashPassword(data.NewPassword)
		if err != nil {
			return ce.NewError(span, ce.CodePasswordHashingFailed, ce.MsgInternalServer, err)
		}

		return u.ar.UpdatePassword(ctx, auth.ID, hashedNewPassword)
	})
}

func (u *authUsecase) ForgotPassword(ctx context.Context, email string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ForgotPassword")
	defer span.End()

	normalizedEmail := utils.NormalizeString(email)
	auth, err := u.ar.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		return err
	}
	if auth.Password == nil {
		// oauth-registered accounts cannot reset password
		// however, attempts at this won't return any error
		return nil
	}

	token := utils.GenerateRandomToken()
	tokenDuration := time.Duration(config.AuthGetResetTokenDuration()) * time.Minute
	if err := u.cache.NewOrReplacePasswordResetToken(ctx, auth.ID, token, tokenDuration); err != nil {
		return err
	}

	return u.email.SendPasswordResetToken(ctx, auth.Email, token)
}

func (u *authUsecase) ResetPassword(ctx context.Context, data *entities.NewPassword) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ResetPassword")
	defer span.End()

	authID, err := u.cache.ConsumePasswordResetToken(ctx, data.Token)
	if err != nil {
		return err
	}

	hashedNewPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return ce.NewError(span, ce.CodePasswordHashingFailed, ce.MsgInternalServer, err)
	}

	return u.ar.UpdatePassword(ctx, authID, hashedNewPassword)
}

func (u *authUsecase) ResendVerification(ctx context.Context, authID int64) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ResendVerification")
	defer span.End()

	auth, err := u.ar.GetByID(ctx, authID)
	if err != nil {
		return err
	}
	if auth.IsVerified {
		return ce.NewError(span, ce.CodeAuthVerified, "Email is already verified.", errors.New("account already verified"))
	}

	token := utils.GenerateRandomToken()
	tokenDuration := time.Duration(config.AuthGetVerifyTokenDuration()) * time.Minute
	if err := u.cache.NewOrReplaceVerificationToken(ctx, auth.ID, token, tokenDuration); err != nil {
		return err
	}

	return u.email.SendVerificationToken(ctx, auth.Email, token)
}

func (u *authUsecase) VerifyEmail(ctx context.Context, token string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "VerifyEmail")
	defer span.End()

	authID, err := u.cache.ConsumeVerificationToken(ctx, token)
	if err != nil {
		return err
	}

	return u.ar.VerifyEmail(ctx, authID)
}

func (u *authUsecase) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "IsEmailRegistered")
	defer span.End()

	normalizedEmail := utils.NormalizeString(email)
	return u.ar.IsEmailRegistered(ctx, normalizedEmail)
}

func (u *authUsecase) IsResetTokenValid(ctx context.Context, token string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "IsResetTokenValid")
	defer span.End()

	tokenKey := utils.CacheCreateKey(constants.CacheKeyPasswordResetToken, token)
	return u.cache.Has(ctx, tokenKey)
}

func (u *authUsecase) RefreshSession(ctx context.Context, sessionToken string) (*entities.AuthToken, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "RefreshSession")
	defer span.End()

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.AuthGetSessionDuration()) * time.Minute

	var token entities.AuthToken
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {
		session, err := u.su.GetSession(ctx, sessionToken)
		if err != nil {
			return err
		}
		if !session.ExpiresAt.After(now) {
			return ce.NewError(span, ce.CodeSessionExpired, ce.MsgInvalidCredentials, ce.ErrSessionExpired)
		}
		if session.RevokedAt != nil {
			return ce.NewError(span, ce.CodeSessionRevoked, ce.MsgInvalidCredentials, ce.ErrSessionRevoked)
		}

		auth, err := u.ar.GetByID(ctx, session.AuthID)
		if err != nil {
			return err
		}

		newSessionToken := utils.GenerateRandomToken()
		newAccessToken, err := utils.GenerateJWTToken(auth.ID, auth.Role, auth.IsVerified)
		if err != nil {
			return ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, err)
		}

		newSessionData := entities.SessionReissue{
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
