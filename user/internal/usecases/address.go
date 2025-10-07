package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/repositories"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/database"
	"go.opentelemetry.io/otel"
)

const addressErrorTracer string = "usecase.address"

type AddressUsecase interface {
	CreateAddress(ctx context.Context, authID int64, data *entities.CreateAddress) (createdAddress *entities.Address, oldPrimaryAddress *entities.Address, err error)
	GetAllAddresses(ctx context.Context, authID int64) (addresses []entities.Address, err error)
	UpdateAddress(ctx context.Context, authID, addressID int64, data *entities.UpdateAddress) (updatedAddress *entities.Address, err error)
	SetPrimaryAddress(ctx context.Context, authID, addressID int64) (newPrimaryAddress *entities.Address, oldPrimaryAddress *entities.Address, err error)
	DeleteAddress(ctx context.Context, authID, addressID int64) (newPrimaryAddress *entities.Address, err error)
}

type addressUsecase struct {
	ar         repositories.AddressRepository
	transactor *database.Transactor
}

func NewAddressUsecase(ar repositories.AddressRepository, transactor *database.Transactor) AddressUsecase {
	return &addressUsecase{ar, transactor}
}

func (u *addressUsecase) CreateAddress(ctx context.Context, authID int64, data *entities.CreateAddress) (*entities.Address, *entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "CreateAddress")
	defer span.End()

	var oldPrimaryAddress *entities.Address
	var address *entities.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) (err error) {
		exists, err := u.ar.HasPrimary(ctx, authID)
		if err != nil {
			return err
		}
		if exists && data.IsPrimary {
			oldPrimaryAddress, err = u.ar.UnsetPrimary(ctx, authID)
			if err != nil {
				return err
			}
		}
		if !exists {
			data.IsPrimary = true
		}

		// TODO (1)

		address, err = u.ar.Create(ctx, authID, data)
		return err
	})

	return address, oldPrimaryAddress, err
}

func (u *addressUsecase) GetAllAddresses(ctx context.Context, authID int64) ([]entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "GetAllAddresses")
	defer span.End()

	return u.ar.GetAll(ctx, authID)
}

func (u *addressUsecase) UpdateAddress(ctx context.Context, authID, addressID int64, data *entities.UpdateAddress) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "UpdateAddress")
	defer span.End()

	// TODO (1)

	return u.ar.Update(ctx, authID, addressID, data)
}

func (u *addressUsecase) SetPrimaryAddress(ctx context.Context, authID, addressID int64) (*entities.Address, *entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "SetPrimaryAddress")
	defer span.End()

	var newPrimaryAddress, oldPrimaryAddress *entities.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) (err error) {
		oldPrimaryAddress, err = u.ar.UnsetPrimary(ctx, authID)
		if err != nil {
			return err
		}

		newPrimaryAddress, err = u.ar.SetPrimary(ctx, authID, addressID)
		return err
	})

	return newPrimaryAddress, oldPrimaryAddress, err
}

func (u *addressUsecase) DeleteAddress(ctx context.Context, authID, addressID int64) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "DeleteAddress")
	defer span.End()

	var newPrimaryAddress *entities.Address
	err := u.transactor.WithTx(ctx, func(ctx context.Context) (err error) {
		if err := u.ar.Delete(ctx, authID, addressID); err != nil {
			return err
		}

		exists, err := u.ar.HasPrimary(ctx, authID)
		if err != nil {
			return err
		}
		if !exists {
			newPrimaryAddress, err = u.ar.SetLastUpdatedPrimary(ctx, authID)
		}

		return err
	})

	return newPrimaryAddress, err
}
