package usecases

import (
	"context"
	"log"
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

const OAuthErrorTracer = ce.OAuthUsecaseTracer

type OAuthUsecase interface {
	Authenticate(ctx context.Context, data *entities.NewOAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
}

type oAuthUsecase struct {
	oar   repos.OAuthRepo
	ar    repos.AuthRepo
	tx    dbtx.TxManager
	su    SessionUsecase
	cache caches.Cache
	email email.EmailService
}

func NewOAuthUsecase(
	oar repos.OAuthRepo,
	ar repos.AuthRepo,
	tx dbtx.TxManager,
	su SessionUsecase,
	cache caches.Cache,
	email email.EmailService,
) OAuthUsecase {
	return &oAuthUsecase{oar, ar, tx, su, cache, email}
}

func (u *oAuthUsecase) Authenticate(ctx context.Context, data *entities.NewOAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	tracer := OAuthErrorTracer + ": Authenticate()"

	if !data.IsVerified {
		return nil, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgEmailNotVerified, tracer, ce.ErrOAuthEmailNotVerified)
	}

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.GetSessionDuration()) * time.Minute
	newAccount := false

	var normalizedEmail string
	var authID int64
	var token entities.AuthToken
	err := u.tx.ReturnError(ctx, func(ctx context.Context) (err error) {
		normalizedEmail = utils.NormalizeString(data.Email)
		exists, auth, err := u.ar.GetForOAuth(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists && auth.Password != nil {
			// Is a regular auth
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidCredentials, tracer, ce.ErrOAuthRegularLogin)
		}
		if exists {
			// Login
			if auth.LockedUntil != nil && auth.LockedUntil.After(now) {
				return ce.NewError(ce.ErrCodeLocked, ce.ErrMsgAccountLocked, tracer, ce.ErrAccountLocked)
			}
		}
		if !exists {
			// Register
			createOAuthData := entities.NewAuth{
				Email: normalizedEmail,
				Role:  constants.RoleCustomer,
			}

			authID, err := u.ar.CreateByOAuth(ctx, &createOAuthData)
			if err != nil {
				return err
			}

			_, err = u.oar.Create(ctx, authID, data)
			if err != nil {
				return err
			}

			newAccount = true
			auth = &entities.Auth{
				ID:         authID,
				IsVerified: false,
				Role:       createOAuthData.Role,
			}
		}

		sessionToken := utils.GenerateRandomToken()
		accessToken, err := utils.GenerateJWTToken(auth.ID, auth.Role, auth.IsVerified)
		if err != nil {
			return ce.NewError(ce.ErrCodeToken, ce.ErrMsgInternalServer, tracer, err)
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
			return err
		}

		authID = auth.ID
		token = entities.AuthToken{
			AccessToken:  accessToken,
			SessionToken: sessionToken,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if newAccount {
		verificationToken := utils.GenerateRandomToken()
		tokenDuration := time.Duration(config.GetEmailVerificationTokenDuration()) * time.Minute
		if err := u.cache.NewOrReplaceVerificationToken(ctx, authID, verificationToken, tokenDuration); err != nil {
			log.Println("WARNING: failed to set verification token in redis after registration:", err.Error())
			return &token, nil
		}

		if err := u.email.SendWelcomeMessage(normalizedEmail, verificationToken); err != nil {
			log.Println("WARNING: failed to send welcome message after registration:", err.Error())
		}
	}

	return &token, nil
}
