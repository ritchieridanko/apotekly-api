package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/apotekly-api/auth/internal/app/repositories"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const sessionErrorTracer string = "usecase.session"

type SessionUsecase interface {
	CreateSession(ctx context.Context, authID int64, data *entities.CreateSession) (err error)
	CreateFirstSession(ctx context.Context, authID int64, data *entities.CreateSession) (err error)
	GetSession(ctx context.Context, token string) (session *entities.Session, err error)
	RevokeSession(ctx context.Context, token string) (err error)
	RefreshSession(ctx context.Context, authID int64, data *entities.CreateSession) (err error)
}

type sessionUsecase struct {
	sr         repositories.SessionRepository
	transactor *database.Transactor
}

func NewSessionUsecase(
	sr repositories.SessionRepository,
	transactor *database.Transactor,
) SessionUsecase {
	return &sessionUsecase{sr, transactor}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, authID int64, data *entities.CreateSession) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "CreateSession")
	defer span.End()

	return u.transactor.WithTx(ctx, func(ctx context.Context) error {
		revokedSessionID, err := u.sr.RevokeActive(ctx, authID)
		if err != nil {
			return err
		}
		if revokedSessionID != 0 {
			data.ParentID = &revokedSessionID
		}

		return u.sr.Create(ctx, authID, data)
	})
}

func (u *sessionUsecase) CreateFirstSession(ctx context.Context, authID int64, data *entities.CreateSession) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "CreateFirstSession")
	defer span.End()

	return u.sr.Create(ctx, authID, data)
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

func (u *sessionUsecase) RefreshSession(ctx context.Context, authID int64, data *entities.CreateSession) error {
	ctx, span := otel.Tracer(sessionErrorTracer).Start(ctx, "RefreshSession")
	defer span.End()

	if data.ParentID == nil {
		err := fmt.Errorf("failed to refresh session: %w", errors.New("session parent id not provided"))
		return ce.NewError(span, ce.CodeSessionNotFound, ce.MsgInvalidCredentials, err)
	}

	return u.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := u.sr.RevokeByID(ctx, *data.ParentID); err != nil {
			return err
		}

		return u.sr.Create(ctx, authID, data)
	})
}
