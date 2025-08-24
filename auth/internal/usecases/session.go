package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const SessionErrorTracer = ce.SessionUsecaseTracer

type SessionUsecase interface {
	NewSessionOnRegister(ctx context.Context, data *entities.NewSession) (err error)
}

type sessionUsecase struct {
	sr repos.SessionRepo
}

func NewSessionUsecase(sr repos.SessionRepo) SessionUsecase {
	return &sessionUsecase{sr}
}

func (u *sessionUsecase) NewSessionOnRegister(ctx context.Context, data *entities.NewSession) error {
	_, err := u.sr.Create(ctx, data)
	return err
}
