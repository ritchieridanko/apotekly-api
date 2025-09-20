package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/db"
	"go.opentelemetry.io/otel"
)

const sessionErrorTracer string = "usecase.session"

type SessionUsecase interface {
	NewSession(ctx context.Context, data *entities.NewSession) (sessionID int64, err error)
	NewSessionOnRegister(ctx context.Context, data *entities.NewSession) (err error)
	RenewSession(ctx context.Context, data *entities.ReissueSession) (newSessionID int64, err error)
	GetSession(ctx context.Context, token string) (session *entities.Session, err error)
	RevokeSession(ctx context.Context, token string) (err error)
}

type sessionUsecase struct {
	sr repos.SessionRepo
	tx db.TxManager
}

func NewSessionUsecase(sr repos.SessionRepo, tx db.TxManager) SessionUsecase {
	return &sessionUsecase{sr, tx}
}

func (u *sessionUsecase) NewSession(ctx context.Context, data *entities.NewSession) (int64, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "NewSession")
	defer span.End()

	var sessionID int64
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {
		revokedSessionID, err := u.sr.RevokeActive(ctx, data.AuthID)
		if err != nil {
			return err
		}
		if revokedSessionID != 0 {
			// there was an active session that's revoked
			reissueData := entities.ReissueSession{
				AuthID:    data.AuthID,
				ParentID:  revokedSessionID,
				Token:     data.Token,
				UserAgent: data.UserAgent,
				IPAddress: data.IPAddress,
				ExpiresAt: data.ExpiresAt,
			}

			sessionID, err = u.sr.Reissue(ctx, &reissueData)
			if err != nil {
				return err
			}

			return nil
		}

		sessionID, err = u.sr.Create(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return sessionID, nil
}

func (u *sessionUsecase) NewSessionOnRegister(ctx context.Context, data *entities.NewSession) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "NewSessionOnRegister")
	defer span.End()

	_, err := u.sr.Create(ctx, data)
	return err
}

func (u *sessionUsecase) RenewSession(ctx context.Context, data *entities.ReissueSession) (int64, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RenewSession")
	defer span.End()

	var newSessionID int64
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {
		err := u.sr.RevokeByID(ctx, data.ParentID)
		if err != nil {
			return err
		}

		newSessionID, err = u.sr.Reissue(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return newSessionID, nil
}

func (u *sessionUsecase) GetSession(ctx context.Context, token string) (*entities.Session, error) {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "GetSession")
	defer span.End()

	return u.sr.GetByToken(ctx, token)
}

func (u *sessionUsecase) RevokeSession(ctx context.Context, token string) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RevokeSession")
	defer span.End()

	return u.sr.RevokeByToken(ctx, token)
}
