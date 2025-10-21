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
)

// TODO
// 1: Send verification token to notification service

const oAuthErrorTracer string = "usecase.oauth"

type OAuthUsecase interface {
	Authenticate(ctx context.Context, data *entities.OAuth, request *entities.Request) (sessionToken string, exchangeCode string, err error)
	ExchangeCode(ctx context.Context, code string) (auth *entities.Auth, accessToken string, err error)
}

type oAuthUsecase struct {
	oar        repositories.OAuthRepository
	ar         repositories.AuthRepository
	oac        caches.OAuthCache
	ac         caches.AuthCache
	su         SessionUsecase
	transactor *database.Transactor
	jwt        *services.JWTService
	cfg        *configs.Config
}

func NewOAuthUsecase(
	oar repositories.OAuthRepository,
	ar repositories.AuthRepository,
	oac caches.OAuthCache,
	ac caches.AuthCache,
	su SessionUsecase,
	transactor *database.Transactor,
	jwt *services.JWTService,
	cfg *configs.Config,
) OAuthUsecase {
	return &oAuthUsecase{oar, ar, oac, ac, su, transactor, jwt, cfg}
}

func (u *oAuthUsecase) Authenticate(ctx context.Context, data *entities.OAuth, request *entities.Request) (string, string, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "Authenticate")
	defer span.End()

	if !data.IsVerified {
		err := fmt.Errorf("failed to authenticate: %w", errors.New("user email not verified"))
		return "", "", ce.NewError(span, ce.CodeOAuthNotVerified, "Cannot authenticate with unverified email", err)
	}

	now := time.Now().UTC()
	newAccount := false

	var rAuth *entities.Auth
	var sessionToken string
	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		normalizedEmail := utils.Normalize(data.Email)
		exists, auth, err := u.ar.GetForOAuth(ctx, normalizedEmail)
		if err != nil {
			return err
		}
		if exists && auth.Password != nil {
			// this is a regular type account
			err := fmt.Errorf("failed to authenticate: %w", errors.New("email registered as regular auth"))
			return ce.NewError(span, ce.CodeOAuthRegularExists, ce.MsgInvalidCredentials, err)
		}
		if !exists {
			// register if not exists
			newAccount = true
			newAuthData := entities.CreateAuth{
				Email:  normalizedEmail,
				RoleID: constants.RoleCustomer,
			}

			auth, err = u.ar.Create(ctx, &newAuthData)
			if err != nil {
				return err
			}
			if err := u.oar.Create(ctx, auth.ID, data); err != nil {
				return err
			}
		}

		sessionToken = utils.NewUUID().String()
		newSessionData := entities.CreateSession{
			Token:     sessionToken,
			UserAgent: request.UserAgent,
			IPAddress: request.IPAddress,
			ExpiresAt: now.Add(u.cfg.Auth.TokenDuration.Session),
		}

		if newAccount {
			err = u.su.CreateFirstSession(ctx, auth.ID, &newSessionData)
		} else {
			err = u.su.CreateSession(ctx, auth.ID, &newSessionData)
		}

		rAuth = auth
		return err
	})
	if err != nil {
		return "", "", err
	}

	exchangeCode := utils.NewUUID().String()
	if err := u.oac.StoreAuth(ctx, exchangeCode, rAuth, u.cfg.OAuth.Duration.CodeExchange); err != nil {
		return "", "", err
	}

	if newAccount {
		verificationToken := utils.NewUUID().String()
		err := u.ac.CreateVerificationToken(
			ctx, rAuth.ID, verificationToken,
			u.cfg.Auth.TokenDuration.Verification,
		)
		if err != nil {
			log.Println("WARNING ->", err.Error())
			return sessionToken, exchangeCode, nil
		}

		// TODO (1)
	}

	return sessionToken, exchangeCode, nil
}

func (u *oAuthUsecase) ExchangeCode(ctx context.Context, code string) (*entities.Auth, string, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "ExchangeCode")
	defer span.End()

	auth, err := u.oac.GetAuth(ctx, code)
	if err != nil {
		return nil, "", err
	}

	accessToken, err := u.jwt.Create(auth.ID, auth.RoleID, auth.IsVerified)
	if err != nil {
		wErr := fmt.Errorf("failed to exchange code: %w", err)
		return nil, "", ce.NewError(span, ce.CodeJWTGenerationFailed, ce.MsgInternalServer, wErr)
	}

	return auth, accessToken, nil
}
