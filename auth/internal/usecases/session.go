package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

type SessionUsecase interface {
	NewSession(ctx context.Context, data *entities.NewSession) (sessionID int64, err error)
	NewSessionOnRegister(ctx context.Context, data *entities.NewSession) (err error)
	RenewSession(ctx context.Context, data *entities.ReissueSession) (newSessionID int64, err error)
	GetSession(ctx context.Context, token string) (session *entities.Session, err error)
	RevokeSession(ctx context.Context, token string) (err error)
}

type sessionUsecase struct {
	sr repos.SessionRepo
	tx dbtx.TxManager
}

func NewSessionUsecase(sr repos.SessionRepo, tx dbtx.TxManager) SessionUsecase {
	return &sessionUsecase{sr, tx}
}

func (u *sessionUsecase) NewSession(ctx context.Context, data *entities.NewSession) (int64, error) {
	var sessionID int64
	err := u.tx.ReturnError(ctx, func(ctx context.Context) error {
		hasAny, activeSessionID, err := u.sr.HasActiveSession(ctx, data.AuthID)
		if err != nil {
			return err
		}
		if hasAny {
			if err := u.sr.RevokeByID(ctx, activeSessionID); err != nil {
				return err
			}

			reissueData := entities.ReissueSession{
				AuthID:    data.AuthID,
				ParentID:  activeSessionID,
				Token:     data.Token,
				UserAgent: data.UserAgent,
				IPAddress: data.IPAddress,
				ExpiresAt: data.ExpiresAt.UTC(),
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
	_, err := u.sr.Create(ctx, data)
	return err
}

func (u *sessionUsecase) RenewSession(ctx context.Context, data *entities.ReissueSession) (int64, error) {
	var newSessionID int64
	err := u.tx.ReturnError(ctx, func(ctx context.Context) error {
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
	return u.sr.GetByToken(ctx, token)
}

func (u *sessionUsecase) RevokeSession(ctx context.Context, token string) error {
	return u.sr.RevokeByToken(ctx, token)
}
