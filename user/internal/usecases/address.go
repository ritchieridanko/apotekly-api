package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"go.opentelemetry.io/otel"
)

// TODO
// 1. Validate the hierarchy of administrative divisions, the postal code, and the coordinates

const addressErrorTracer string = "usecase.address"

type AddressUsecase interface {
	NewAddress(ctx context.Context, authID int64, data *entities.NewAddress) (address *entities.Address, err error)
}

type addressUsecase struct {
	ar repos.AddressRepo
	tx db.TxManager
}

func NewAddressUsecase(ar repos.AddressRepo, tx db.TxManager) AddressUsecase {
	return &addressUsecase{ar, tx}
}

func (u *addressUsecase) NewAddress(ctx context.Context, authID int64, data *entities.NewAddress) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "NewAddress")
	defer span.End()

	var address *entities.Address
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		hasPrimaryAddress, err := u.ar.HasPrimaryAddress(ctx, authID)
		if err != nil {
			return err
		}
		if hasPrimaryAddress && data.IsPrimary {
			if err := u.ar.UnsetPrimaryAddress(ctx, authID); err != nil {
				return err
			}
		}
		if !hasPrimaryAddress {
			data.IsPrimary = true
		}

		// TODO (1)

		address, err = u.ar.Create(ctx, authID, data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return address, nil
}
