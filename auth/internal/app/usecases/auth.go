package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/caches"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/repositories"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TODO
// 1: Send verification token to notification service
// 2: Send email change token to notification service
// 3: Send reset token to notification service

const authErrorTracer string = "usecase.auth"

type AuthUsecase interface {
	Register(ctx context.Context, data *entities.CreateAuth, request *entities.Request) (authToken *entities.AuthToken, createdAuth *entities.Auth, err error)
	Login(ctx context.Context, data *entities.GetAuth, request *entities.Request) (authToken *entities.AuthToken, auth *entities.Auth, err error)
	Logout(ctx context.Context, sessionToken string) (err error)
	ChangeEmail(ctx context.Context, authID int64, email string) (recipientEmail string, err error)
	ConfirmEmailChange(ctx context.Context, token, sessionToken string) (authToken *entities.AuthToken, updatedAuth *entities.Auth, err error)
	ChangePassword(ctx context.Context, authID int64, data *entities.UpdatePassword) (err error)
	ForgotPassword(ctx context.Context, email string) (recipientEmail string, err error)
	ResetPassword(ctx context.Context, data *entities.ResetPassword) (err error)
	ResendVerification(ctx context.Context, authID int64) (recipientEmail string, err error)
	VerifyAccount(ctx context.Context, token, sessionToken string) (authToken *entities.AuthToken, updatedAuth *entities.Auth, err error)
	RefreshSession(ctx context.Context, sessionToken string) (authToken *entities.AuthToken, err error)
	IsEmailRegistered(ctx context.Context, email string) (isRegistered bool, err error)
	IsResetTokenValid(ctx context.Context, token string) (isValid bool, err error)
}

type authUsecase struct {
	ar         repositories.AuthRepository
	ac         caches.AuthCache
	su         SessionUsecase
	transactor *database.Transactor
	bcrypt     *services.BCryptService
	jwt        *services.JWTService
	cfg        *configs.Config
}

func NewAuthUsecase(
	ar repositories.AuthRepository,
	ac caches.AuthCache,
	su SessionUsecase,
	transactor *database.Transactor,
	bcrypt *services.BCryptService,
	jwt *services.JWTService,
	cfg *configs.Config,
) AuthUsecase {
	return &authUsecase{ar, ac, su, transactor, bcrypt, jwt, cfg}
}

func (u *authUsecase) Register(ctx context.Context, data *entities.CreateAuth, request *entities.Request) (*entities.AuthToken, *entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Register")
	defer span.End()

	now := time.Now().UTC()

	var auth *entities.Auth
	var authToken entities.AuthToken
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		normalizedEmail := utils.Normalize(data.Email)
		exists, err := u.ar.Exists(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists {
			err := fmt.Errorf("failed to register: %w", errors.New("email conflict"))
			return ce.NewError(span, ce.CodeAuthEmailConflict, "Email is already registered", err)
		}

		hashedPassword, err := u.bcrypt.Hash(*data.Password)
		if err != nil {
			wErr := fmt.Errorf("failed to register: %w", err)
			return ce.NewError(span, ce.CodePasswordHashingFailed, ce.MsgInternalServer, wErr)
		}

		newAuthData := entities.CreateAuth{
			Email:    normalizedEmail,
			Password: &hashedPassword,
			RoleID:   constants.RoleCustomer,
		}

		auth, err = u.ar.Create(ctx, &newAuthData)
		if err != nil {
			return err
		}

		sessionToken := utils.NewUUID().String()
		accessToken, err := u.jwt.Create(auth.ID, auth.RoleID, auth.IsVerified)
		if err != nil {
			wErr := fmt.Errorf("failed to register: %w", err)
			return ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, wErr)
		}

		newSessionData := entities.CreateSession{
			Token:     sessionToken,
			UserAgent: request.UserAgent,
			IPAddress: request.IPAddress,
			ExpiresAt: now.Add(u.cfg.Auth.TokenDuration.Session),
		}
		if err := u.su.CreateFirstSession(ctx, auth.ID, &newSessionData); err != nil {
			return err
		}

		authToken = entities.AuthToken{
			AccessToken:  accessToken,
			SessionToken: sessionToken,
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	verificationToken := utils.NewUUID().String()
	err = u.ac.CreateVerificationToken(
		ctx, auth.ID, verificationToken,
		u.cfg.Auth.TokenDuration.Verification,
	)
	if err != nil {
		log.Println("WARNING -> ", err.Error())
		return &authToken, auth, nil
	}

	// TODO (1)

	return &authToken, auth, nil
}

func (u *authUsecase) Login(ctx context.Context, data *entities.GetAuth, request *entities.Request) (*entities.AuthToken, *entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Login")
	defer span.End()

	now := time.Now().UTC()

	normalizedEmail := utils.Normalize(data.Email)
	auth, err := u.ar.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		return nil, nil, err
	}
	if auth.Password == nil {
		// this is an oauth type account
		err := fmt.Errorf("failed to login: %w", errors.New("email registered as oauth"))
		return nil, nil, ce.NewError(span, ce.CodeOAuthRegularLogin, ce.MsgInvalidCredentials, err)
	}
	if err := u.bcrypt.Validate(*auth.Password, data.Password); err != nil {
		wErr := fmt.Errorf("failed to login: %w", err)
		return nil, nil, ce.NewError(span, ce.CodeAuthWrongPassword, ce.MsgInvalidCredentials, wErr)
	}

	sessionToken := utils.NewUUID().String()
	accessToken, err := u.jwt.Create(auth.ID, auth.RoleID, auth.IsVerified)
	if err != nil {
		wErr := fmt.Errorf("failed to login: %w", err)
		return nil, nil, ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, wErr)
	}

	newSessionData := entities.CreateSession{
		Token:     sessionToken,
		UserAgent: request.UserAgent,
		IPAddress: request.IPAddress,
		ExpiresAt: now.Add(u.cfg.Auth.TokenDuration.Session),
	}
	if err := u.su.CreateSession(ctx, auth.ID, &newSessionData); err != nil {
		return nil, nil, err
	}

	authToken := entities.AuthToken{
		AccessToken:  accessToken,
		SessionToken: sessionToken,
	}

	return &authToken, auth, nil
}

func (u *authUsecase) Logout(ctx context.Context, sessionToken string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "Logout")
	defer span.End()

	return u.su.RevokeSession(ctx, sessionToken)
}

func (u *authUsecase) ChangeEmail(ctx context.Context, authID int64, email string) (string, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ChangeEmail")
	defer span.End()

	auth, err := u.ar.GetByID(ctx, authID)
	if err != nil {
		return "", err
	}
	if auth.Password == nil {
		// this is an oauth type account
		// email cannot be changed for oauth accounts
		err := fmt.Errorf("failed to change email: %w", errors.New("email change with oauth account"))
		return "", ce.NewError(span, ce.CodeOAuthEmailChange, "OAuth account cannot change email", err)
	}

	normalizedEmail := utils.Normalize(email)
	exists, err := u.ar.Exists(ctx, normalizedEmail)
	if err != nil {
		return "", err
	}
	if exists {
		err := fmt.Errorf("failed to change email: %w", errors.New("email conflict"))
		return "", ce.NewError(span, ce.CodeAuthEmailConflict, "Email is already registered", err)
	}

	token := utils.NewUUID().String()
	err = u.ac.CreateEmailChangeToken(
		ctx, auth.ID, normalizedEmail, token,
		u.cfg.Auth.TokenDuration.EmailChange,
	)
	if err != nil {
		return "", err
	}

	// TODO (2)

	return normalizedEmail, nil
}

func (u *authUsecase) ConfirmEmailChange(ctx context.Context, token, sessionToken string) (*entities.AuthToken, *entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ConfirmEmailChange")
	defer span.End()

	authID, newEmail, err := u.ac.UseEmailChangeToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	var auth *entities.Auth
	var authToken *entities.AuthToken
	err = u.transactor.WithTx(ctx, func(ctx context.Context) error {
		auth, err = u.ar.UpdateEmail(ctx, authID, newEmail)
		if err != nil {
			return err
		}

		if sessionToken != "" {
			authToken, err = u.RefreshSession(ctx, sessionToken)
			if err != nil {
				// non-fatal: trace the failure, but continue
				span.AddEvent(
					"RefreshSession failed, continuing without a new session token",
					trace.WithAttributes(attribute.String("error", err.Error())),
				)
			}
		}
		return nil
	})

	return authToken, auth, err
}

func (u *authUsecase) ChangePassword(ctx context.Context, authID int64, data *entities.UpdatePassword) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ChangePassword")
	defer span.End()

	return u.transactor.WithTx(ctx, func(ctx context.Context) error {
		auth, err := u.ar.GetByID(ctx, authID)
		if err != nil {
			return err
		}
		if auth.Password == nil {
			// this is an oauth type account
			// password cannot be changed for oauth accounts
			err := fmt.Errorf("failed to change password: %w", errors.New("password change with oauth account"))
			return ce.NewError(span, ce.CodeOAuthPasswordChange, "OAuth account cannot change password", err)
		}
		if err := u.bcrypt.Validate(*auth.Password, data.OldPassword); err != nil {
			wErr := fmt.Errorf("failed to change password: %w", err)
			return ce.NewError(span, ce.CodeAuthWrongPassword, "Invalid old password", wErr)
		}

		hashedNewPassword, err := u.bcrypt.Hash(data.NewPassword)
		if err != nil {
			wErr := fmt.Errorf("failed to change password: %w", err)
			return ce.NewError(span, ce.CodePasswordHashingFailed, ce.MsgInternalServer, wErr)
		}

		return u.ar.UpdatePassword(ctx, auth.ID, hashedNewPassword)
	})
}

func (u *authUsecase) ForgotPassword(ctx context.Context, email string) (string, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ForgotPassword")
	defer span.End()

	normalizedEmail := utils.Normalize(email)
	auth, err := u.ar.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		return "", err
	}
	if auth.Password == nil {
		// this is an oauth type account
		// password cannot be changed for oauth accounts
		// however, attempts at this won't return any error
		return "", nil
	}

	token := utils.NewUUID().String()
	if err := u.ac.CreateResetToken(ctx, auth.ID, token, u.cfg.Auth.TokenDuration.Reset); err != nil {
		return "", err
	}

	// TODO (3)

	return normalizedEmail, nil
}

func (u *authUsecase) ResetPassword(ctx context.Context, data *entities.ResetPassword) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ResetPassword")
	defer span.End()

	authID, err := u.ac.UseResetToken(ctx, data.Token)
	if err != nil {
		return err
	}

	hashedNewPassword, err := u.bcrypt.Hash(data.NewPassword)
	if err != nil {
		wErr := fmt.Errorf("failed to reset password: %w", err)
		return ce.NewError(span, ce.CodePasswordHashingFailed, ce.MsgInternalServer, wErr)
	}

	return u.ar.UpdatePassword(ctx, authID, hashedNewPassword)
}

func (u *authUsecase) ResendVerification(ctx context.Context, authID int64) (string, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ResendVerification")
	defer span.End()

	auth, err := u.ar.GetByID(ctx, authID)
	if err != nil {
		return "", err
	}
	if auth.IsVerified {
		err := fmt.Errorf("failed to resend verification: %w", errors.New("account already verified"))
		return "", ce.NewError(span, ce.CodeAuthVerified, "Email is already verified", err)
	}

	token := utils.NewUUID().String()
	err = u.ac.CreateVerificationToken(
		ctx, auth.ID, token,
		u.cfg.Auth.TokenDuration.Verification,
	)
	if err != nil {
		return "", err
	}

	// TODO (1)

	return auth.Email, nil
}

func (u *authUsecase) VerifyAccount(ctx context.Context, token, sessionToken string) (*entities.AuthToken, *entities.Auth, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "VerifyAccount")
	defer span.End()

	authID, err := u.ac.UseVerificationToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	var auth *entities.Auth
	var authToken *entities.AuthToken
	err = u.transactor.WithTx(ctx, func(ctx context.Context) error {
		auth, err = u.ar.SetVerified(ctx, authID)
		if err != nil {
			return err
		}

		if sessionToken != "" {
			authToken, err = u.RefreshSession(ctx, sessionToken)
			if err != nil {
				// non-fatal: trace the failure, but continue
				span.AddEvent(
					"RefreshSession failed, continuing without a new session token",
					trace.WithAttributes(attribute.String("error", err.Error())),
				)
			}
		}
		return nil
	})

	return authToken, auth, err
}

func (u *authUsecase) RefreshSession(ctx context.Context, sessionToken string) (*entities.AuthToken, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "RefreshSession")
	defer span.End()

	now := time.Now().UTC()

	var authToken entities.AuthToken
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		session, err := u.su.GetSession(ctx, sessionToken)
		if err != nil {
			return err
		}
		if !session.ExpiresAt.After(now) {
			err := fmt.Errorf("failed to refresh session: %w", ce.ErrSessionExpired)
			return ce.NewError(span, ce.CodeSessionExpired, ce.MsgUnauthenticated, err)
		}
		if session.RevokedAt != nil {
			err := fmt.Errorf("failed to refresh session: %w", ce.ErrSessionRevoked)
			return ce.NewError(span, ce.CodeSessionRevoked, ce.MsgUnauthenticated, err)
		}

		auth, err := u.ar.GetByID(ctx, session.AuthID)
		if err != nil {
			return err
		}

		newSessionToken := utils.NewUUID().String()
		newAccessToken, err := u.jwt.Create(auth.ID, auth.RoleID, auth.IsVerified)
		if err != nil {
			wErr := fmt.Errorf("failed to refresh session: %w", err)
			return ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, wErr)
		}

		newSessionData := entities.CreateSession{
			ParentID:  &session.ID,
			Token:     newSessionToken,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			ExpiresAt: now.Add(u.cfg.Auth.TokenDuration.Session),
		}
		if err := u.su.RefreshSession(ctx, auth.ID, &newSessionData); err != nil {
			return err
		}

		authToken = entities.AuthToken{
			AccessToken:  newAccessToken,
			SessionToken: newSessionToken,
		}

		return nil
	})

	return &authToken, err
}

func (u *authUsecase) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "IsEmailRegistered")
	defer span.End()

	normalizedEmail := utils.Normalize(email)
	return u.ar.Exists(ctx, normalizedEmail)
}

func (u *authUsecase) IsResetTokenValid(ctx context.Context, token string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "IsResetTokenValid")
	defer span.End()

	return u.ac.ResetTokenExists(ctx, token)
}
