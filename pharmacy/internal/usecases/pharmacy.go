package usecases

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"

	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/ce"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/entities"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/repos"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/services/storage"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/utils"
	"go.opentelemetry.io/otel"
)

// TODO
// 1. Validate the hierarchy of administrative divisions, the postal code, and the coordinates

const pharmacyErrorTracer string = "usecase.pharmacy"

type PharmacyUsecase interface {
	NewPharmacy(ctx context.Context, authID int64, data *entities.NewPharmacy, image multipart.File) (pharmacy *entities.Pharmacy, err error)
	GetPharmacy(ctx context.Context, authID int64) (pharmacy *entities.Pharmacy, err error)
	UpdatePharmacy(ctx context.Context, authID int64, data *entities.PharmacyChange) (pharmacy *entities.Pharmacy, err error)
	ChangeLogo(ctx context.Context, authID int64, image multipart.File) (err error)
}

type pharmacyUsecase struct {
	pr      repos.PharmacyRepo
	tx      db.TxManager
	storage storage.StorageService
}

func NewPharmacyUsecase(pr repos.PharmacyRepo, tx db.TxManager, storage storage.StorageService) PharmacyUsecase {
	return &pharmacyUsecase{pr, tx, storage}
}

func (u *pharmacyUsecase) NewPharmacy(ctx context.Context, authID int64, data *entities.NewPharmacy, image multipart.File) (*entities.Pharmacy, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "NewPharmacy")
	defer span.End()

	var pharmacy *entities.Pharmacy
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		exists, err := u.pr.HasPharmacy(ctx, authID)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(span, ce.CodeDBDuplicateData, "Pharmacy already exists.", errors.New("auth id conflict"))
		}

		// upload image if exists
		var pictureURL *string
		pharmacyPublicID := utils.GenerateRandomUUID()
		if image != nil {
			imageURL, err := u.uploadImage(ctx, image, pharmacyPublicID.String(), "logo", "pharmacies/logos", true)
			if err != nil {
				log.Println("WARNING -> failed to upload image:", err.Error())
			} else {
				pictureURL = &imageURL
			}
		}

		// TODO (1)

		newData := entities.NewPharmacy{
			PharmacyPublicID: pharmacyPublicID,
			Name:             data.Name,
			LegalName:        data.LegalName,
			Description:      data.Description,
			LicenseNumber:    data.LicenseNumber,
			LicenseAuthority: data.LicenseAuthority,
			LicenseExpiry:    data.LicenseExpiry,
			Email:            data.Email,
			Phone:            data.Phone,
			Website:          data.Website,
			Country:          data.Country,
			AdminLevel1:      data.AdminLevel1,
			AdminLevel2:      data.AdminLevel2,
			AdminLevel3:      data.AdminLevel3,
			AdminLevel4:      data.AdminLevel4,
			Street:           data.Street,
			PostalCode:       data.PostalCode,
			Latitude:         data.Latitude,
			Longitude:        data.Longitude,
			Logo:             pictureURL,
			OpeningHours:     data.OpeningHours,
		}

		pharmacy, err = u.pr.Create(ctx, authID, &newData)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pharmacy, nil
}

func (u *pharmacyUsecase) GetPharmacy(ctx context.Context, authID int64) (*entities.Pharmacy, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "GetPharmacy")
	defer span.End()

	return u.pr.GetByAuthID(ctx, authID)
}

func (u *pharmacyUsecase) UpdatePharmacy(ctx context.Context, authID int64, data *entities.PharmacyChange) (*entities.Pharmacy, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "UpdatePharmacy")
	defer span.End()

	// TODO (1)

	return u.pr.UpdatePharmacy(ctx, authID, data)
}

func (u *pharmacyUsecase) ChangeLogo(ctx context.Context, authID int64, image multipart.File) error {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "ChangeLogo")
	defer span.End()

	publicID, err := u.pr.GetPublicID(ctx, authID)
	if err != nil {
		return err
	}

	imageURL, err := u.uploadImage(ctx, image, publicID.String(), "logo", "pharmacies/logos", true)
	if err != nil {
		return err
	}

	return u.pr.UpdateLogo(ctx, authID, imageURL)
}

func (u *pharmacyUsecase) uploadImage(ctx context.Context, image multipart.File, publicID, prefix, folder string, overwrite bool) (imageURL string, err error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "uploadImage")
	defer span.End()

	imageBuf, err := io.ReadAll(image)
	if err != nil {
		return "", ce.NewError(span, ce.CodeFileBuffer, ce.MsgInternalServer, err)
	}

	if err := utils.ValidateImageFile(imageBuf); err != nil {
		return "", ce.NewError(span, ce.CodeInvalidPayload, "Invalid file type.", err)
	}

	data := entities.NewUpload{
		File:           bytes.NewReader(imageBuf),
		PublicID:       publicID,
		PublicIDPrefix: prefix,
		Folder:         folder,
		Overwrite:      &overwrite,
	}

	result, err := u.storage.Upload(ctx, &data)
	if err != nil {
		return "", ce.NewError(span, ce.CodeFileUploadFailed, ce.MsgInternalServer, err)
	}

	return result.SecureURL, nil
}
