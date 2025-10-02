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
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
)

func (u *authUsecase) RegisterAsPharmacy(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "RegisterAsPharmacy")
	defer span.End()

	now := time.Now().UTC()
	sessionDuration := time.Duration(config.AuthGetSessionDuration()) * time.Minute

	var normalizedEmail string
	var authID int64
	var token entities.AuthToken
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {
		normalizedEmail = utils.Normalize(data.Email)
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
			Role:     constants.RolePharmacy,
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
	if err := u.cache.NewVerificationToken(ctx, authID, verificationToken, tokenDuration); err != nil {
		log.Println("WARNING -> failed to store verification token in cache after registration:", err.Error())
		return &token, nil
	}

	if err := u.email.SendWelcomeMessageForPharmacy(ctx, normalizedEmail, verificationToken); err != nil {
		log.Println("WARNING -> failed to send welcome email after registration:", err.Error())
	}

	return &token, nil
}
