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

const oAuthErrorTracer string = "usecase.oauth"

type OAuthUsecase interface {
	Authenticate(ctx context.Context, data *entities.OAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
}

type oAuthUsecase struct {
	oar   repos.OAuthRepo
	ar    repos.AuthRepo
	su    SessionUsecase
	tx    db.TxManager
	cache cache.CacheService
	email email.EmailService
}

func NewOAuthUsecase(
	oar repos.OAuthRepo,
	ar repos.AuthRepo,
	su SessionUsecase,
	tx db.TxManager,
	cache cache.CacheService,
	email email.EmailService,
) OAuthUsecase {
	return &oAuthUsecase{oar, ar, su, tx, cache, email}
}

func (u *oAuthUsecase) Authenticate(ctx context.Context, data *entities.OAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "Authenticate")
	defer span.End()

	if !data.IsVerified {
		return nil, ce.NewError(span, ce.CodeOAuthNotVerified, "Your email is not verified yet.", errors.New("user email not verified"))
	}

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.AuthGetSessionDuration()) * time.Minute
	newAccount := false

	var normalizedEmail string
	var authID int64
	var token entities.AuthToken
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		normalizedEmail = utils.Normalize(data.Email)
		exists, auth, err := u.ar.GetForOAuth(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists && auth.Password != nil {
			// this is a regular type auth
			return ce.NewError(span, ce.CodeOAuthRegularExists, ce.MsgInvalidCredentials, errors.New("email registered as regular auth"))
		}
		if !exists {
			// register if not exists
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
			return ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, err)
		}

		sessionData := entities.NewSession{
			AuthID:    auth.ID,
			Token:     sessionToken,
			UserAgent: request.UserAgent,
			IPAddress: request.IPAddress,
			ExpiresAt: now.Add(sessionDuration),
		}

		if newAccount {
			err = u.su.NewSessionOnRegister(ctx, &sessionData)
		} else {
			_, err = u.su.NewSession(ctx, &sessionData)
		}
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
		tokenDuration := time.Duration(config.AuthGetVerifyTokenDuration()) * time.Minute
		if err := u.cache.NewVerificationToken(ctx, authID, verificationToken, tokenDuration); err != nil {
			log.Println("WARNING -> failed to store verification token in cache after registration:", err.Error())
			return &token, nil
		}

		if err := u.email.SendWelcomeMessage(ctx, normalizedEmail, verificationToken); err != nil {
			log.Println("WARNING -> failed to send welcome email after registration:", err.Error())
		}
	}

	return &token, nil
}
