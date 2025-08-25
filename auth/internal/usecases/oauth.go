package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
)

type OAuthUsecase interface {
	IsAuthRegistered(ctx context.Context, authID int64) (exists bool, err error)
}

type oAuthUsecase struct {
	oar repos.OAuthRepo
}

func NewOAuthUsecase(oar repos.OAuthRepo) OAuthUsecase {
	return &oAuthUsecase{oar}
}

func (u *oAuthUsecase) IsAuthRegistered(ctx context.Context, authID int64) (bool, error) {
	return u.oar.IsAuthRegistered(ctx, authID)
}
