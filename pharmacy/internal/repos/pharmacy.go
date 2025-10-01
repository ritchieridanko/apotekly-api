package repos

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/ce"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/entities"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/services/db"
	"go.opentelemetry.io/otel"
)

const pharmacyErrorTracer string = "repo.pharmacy"

type PharmacyRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewPharmacy) (pharmacy *entities.Pharmacy, err error)
	HasPharmacy(ctx context.Context, authID int64) (exists bool, err error)
}

type pharmacyRepo struct {
	database db.DBService
}

func NewPharmacyRepo(database db.DBService) PharmacyRepo {
	return &pharmacyRepo{database}
}

func (r *pharmacyRepo) Create(ctx context.Context, authID int64, data *entities.NewPharmacy) (*entities.Pharmacy, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO
			pharmacies (
				auth_id, pharmacy_public_id, name, legal_name, description, license_number,
				license_authority, license_expiry, email, phone, website, country,
				admin_level_1, admin_level_2, admin_level_3, admin_level_4, street,
				postal_code, latitude, longitude, location, logo, opening_hours
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, ST_SetSRID(ST_MakePoint($20, $19), 4326), $21, $22
		)
		RETURNING
			pharmacy_public_id, name, legal_name, description, license_number,
			license_authority, license_expiry, email, phone, website, country,
			admin_level_1, admin_level_2, admin_level_3, admin_level_4, street,
			postal_code, latitude, longitude, logo, opening_hours, status
	`

	row := r.database.QueryRow(
		ctx, query,
		authID, data.PharmacyPublicID, data.Name, data.LegalName, data.Description, data.LicenseNumber,
		data.LicenseAuthority, data.LicenseExpiry, data.Email, data.Phone, data.Website, data.Country,
		data.AdminLevel1, data.AdminLevel2, data.AdminLevel3, data.AdminLevel4, data.Street, data.PostalCode,
		data.Latitude, data.Longitude, data.Logo, data.OpeningHours,
	)

	var pharmacy entities.Pharmacy
	err := row.Scan(
		&pharmacy.PharmacyPublicID, &pharmacy.Name, &pharmacy.LegalName, &pharmacy.Description, &pharmacy.LicenseNumber,
		&pharmacy.LicenseAuthority, &pharmacy.LicenseExpiry, &pharmacy.Email, &pharmacy.Phone, &pharmacy.Website,
		&pharmacy.Country, &pharmacy.AdminLevel1, &pharmacy.AdminLevel2, &pharmacy.AdminLevel3, &pharmacy.AdminLevel4,
		&pharmacy.Street, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude, &pharmacy.Logo,
		&pharmacy.OpeningHours, &pharmacy.Status,
	)
	if err != nil {
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &pharmacy, nil
}

func (r *pharmacyRepo) HasPharmacy(ctx context.Context, authID int64) (bool, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "HasPharmacy")
	defer span.End()

	query := "SELECT 1 FROM pharmacies WHERE auth_id = $1 AND deleted_at IS NULL"
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return true, nil
}
