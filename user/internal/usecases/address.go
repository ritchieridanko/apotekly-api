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
	NewAddress(ctx context.Context, authID int64, data *entities.NewAddress) (address *entities.Address, unsetPrimaryID int64, err error)
	GetAllAddresses(ctx context.Context, authID int64) (addresses []entities.Address, err error)
	UpdateAddress(ctx context.Context, authID, addressID int64, data *entities.AddressChange) (address *entities.Address, unsetPrimaryID int64, err error)
	DeleteAddress(ctx context.Context, authID, addressID int64) (deletedID int64, newPrimaryID int64, err error)
}

type addressUsecase struct {
	ar repos.AddressRepo
	tx db.TxManager
}

func NewAddressUsecase(ar repos.AddressRepo, tx db.TxManager) AddressUsecase {
	return &addressUsecase{ar, tx}
}

func (u *addressUsecase) NewAddress(ctx context.Context, authID int64, data *entities.NewAddress) (*entities.Address, int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "NewAddress")
	defer span.End()

	var unsetPrimaryID int64
	var address *entities.Address
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		hasPrimaryAddress, err := u.ar.HasPrimary(ctx, authID)
		if err != nil {
			return err
		}
		if hasPrimaryAddress && data.IsPrimary {
			unsetPrimaryID, err = u.ar.UnsetPrimary(ctx, authID)
			if err != nil {
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
		return nil, 0, err
	}

	return address, unsetPrimaryID, nil
}

func (u *addressUsecase) GetAllAddresses(ctx context.Context, authID int64) ([]entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "GetAllAddresses")
	defer span.End()

	return u.ar.GetAll(ctx, authID)
}

func (u *addressUsecase) UpdateAddress(ctx context.Context, authID, addressID int64, data *entities.AddressChange) (*entities.Address, int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "UpdateAddress")
	defer span.End()

	var unsetPrimaryID int64
	var address *entities.Address
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		if data.IsPrimary != nil && *data.IsPrimary {
			hasPrimaryAddress, err := u.ar.HasPrimary(ctx, authID)
			if err != nil {
				return err
			}
			if hasPrimaryAddress {
				unsetPrimaryID, err = u.ar.UnsetPrimary(ctx, authID)
				if err != nil {
					return err
				}
			}
		}

		// TODO (1)

		address, err = u.ar.Update(ctx, authID, addressID, data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return address, unsetPrimaryID, nil
}

func (u *addressUsecase) DeleteAddress(ctx context.Context, authID, addressID int64) (int64, int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "DeleteAddress")
	defer span.End()

	var deletedID, newPrimaryID int64
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		deletedID, err = u.ar.Delete(ctx, authID, addressID)
		if err != nil {
			return err
		}

		hasPrimaryAddress, err := u.ar.HasPrimary(ctx, authID)
		if err != nil {
			return err
		}
		if !hasPrimaryAddress {
			newPrimaryID, err = u.ar.SetLastAsPrimary(ctx, authID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return deletedID, newPrimaryID, nil
}
