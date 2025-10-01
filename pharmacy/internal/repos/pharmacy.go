package repos

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/ce"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/entities"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/services/db"
	"go.opentelemetry.io/otel"
)

const pharmacyErrorTracer string = "repo.pharmacy"

type PharmacyRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewPharmacy) (pharmacy *entities.Pharmacy, err error)
	GetByAuthID(ctx context.Context, authID int64) (pharmacy *entities.Pharmacy, err error)
	GetPublicID(ctx context.Context, authID int64) (publicID uuid.UUID, err error)
	UpdatePharmacy(ctx context.Context, authID int64, data *entities.PharmacyChange) (pharmacy *entities.Pharmacy, err error)
	UpdateLogo(ctx context.Context, authID int64, logo string) (err error)
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

func (r *pharmacyRepo) GetByAuthID(ctx context.Context, authID int64) (*entities.Pharmacy, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "GetByAuthID")
	defer span.End()

	query := `
		SELECT
			pharmacy_public_id, name, legal_name, description, license_number,
			license_authority, license_expiry, email, phone, website, country,
			admin_level_1, admin_level_2, admin_level_3, admin_level_4, street,
			postal_code, latitude, longitude, logo, opening_hours, status
		FROM
			pharmacies
		WHERE
			auth_id = $1
			AND deleted_at IS NULL
	`
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var pharmacy entities.Pharmacy
	err := row.Scan(
		&pharmacy.PharmacyPublicID, &pharmacy.Name, &pharmacy.LegalName, &pharmacy.Description, &pharmacy.LicenseNumber,
		&pharmacy.LicenseAuthority, &pharmacy.LicenseExpiry, &pharmacy.Email, &pharmacy.Phone, &pharmacy.Website,
		&pharmacy.Country, &pharmacy.AdminLevel1, &pharmacy.AdminLevel2, &pharmacy.AdminLevel3, &pharmacy.AdminLevel4,
		&pharmacy.Street, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude, &pharmacy.Logo,
		&pharmacy.OpeningHours, &pharmacy.Status,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodePharmacyNotFound, ce.MsgPharmacyNotFound, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &pharmacy, nil
}

func (r *pharmacyRepo) GetPublicID(ctx context.Context, authID int64) (uuid.UUID, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "GetPublicID")
	defer span.End()

	query := "SELECT pharmacy_public_id FROM pharmacies WHERE auth_id = $1 AND deleted_at IS NULL"
	if r.database.IsWithinTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var publicID uuid.UUID
	if err := row.Scan(&publicID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return uuid.Nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return uuid.Nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return publicID, nil
}

func (r *pharmacyRepo) UpdatePharmacy(ctx context.Context, authID int64, data *entities.PharmacyChange) (*entities.Pharmacy, error) {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "UpdatePharmacy")
	defer span.End()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if data.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *data.Name)
		argPos++
	}
	if data.LegalName != nil {
		setClauses = append(setClauses, fmt.Sprintf("legal_name = $%d", argPos))
		args = append(args, *data.LegalName)
		argPos++
	}
	if data.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *data.Description)
		argPos++
	}
	if data.LicenseNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("license_number = $%d", argPos))
		args = append(args, *data.LicenseNumber)
		argPos++
	}
	if data.LicenseAuthority != nil {
		setClauses = append(setClauses, fmt.Sprintf("license_authority = $%d", argPos))
		args = append(args, *data.LicenseAuthority)
		argPos++
	}
	if data.LicenseExpiry != nil {
		setClauses = append(setClauses, fmt.Sprintf("license_expiry = $%d", argPos))
		args = append(args, *data.LicenseExpiry)
		argPos++
	}
	if data.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argPos))
		args = append(args, *data.Email)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if data.Website != nil {
		setClauses = append(setClauses, fmt.Sprintf("website = $%d", argPos))
		args = append(args, *data.Website)
		argPos++
	}
	if data.Country != nil {
		setClauses = append(setClauses, fmt.Sprintf("country = $%d", argPos))
		args = append(args, *data.Country)
		argPos++
	}
	if data.AdminLevel1 != nil {
		setClauses = append(setClauses, fmt.Sprintf("admin_level_1 = $%d", argPos))
		args = append(args, *data.AdminLevel1)
		argPos++
	}
	if data.AdminLevel2 != nil {
		setClauses = append(setClauses, fmt.Sprintf("admin_level_2 = $%d", argPos))
		args = append(args, *data.AdminLevel2)
		argPos++
	}
	if data.AdminLevel3 != nil {
		setClauses = append(setClauses, fmt.Sprintf("admin_level_3 = $%d", argPos))
		args = append(args, *data.AdminLevel3)
		argPos++
	}
	if data.AdminLevel4 != nil {
		setClauses = append(setClauses, fmt.Sprintf("admin_level_4 = $%d", argPos))
		args = append(args, *data.AdminLevel4)
		argPos++
	}
	if data.Street != nil {
		setClauses = append(setClauses, fmt.Sprintf("street = $%d", argPos))
		args = append(args, *data.Street)
		argPos++
	}
	if data.PostalCode != nil {
		setClauses = append(setClauses, fmt.Sprintf("postal_code = $%d", argPos))
		args = append(args, *data.PostalCode)
		argPos++
	}
	if data.Latitude != nil && data.Longitude != nil {
		setClauses = append(setClauses,
			fmt.Sprintf("latitude = $%d", argPos),
			fmt.Sprintf("longitude = $%d", argPos+1),
			fmt.Sprintf("location = ST_SetSRID(ST_MakePoint($%d, $%d), 4326)", argPos+1, argPos),
		)
		args = append(args, *data.Latitude, *data.Longitude)
		argPos += 2
	}
	if data.OpeningHours != nil {
		setClauses = append(setClauses, fmt.Sprintf("opening_hours = $%d", argPos))
		args = append(args, *data.OpeningHours)
		argPos++
	}
	if len(setClauses) == 0 {
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgNoFieldsToUpdate, ce.ErrNoFieldsProvided)
	}
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
			UPDATE pharmacies
			SET %s
			WHERE auth_id = $%d AND deleted_at IS NULL
			RETURNING
				pharmacy_public_id, name, legal_name, description, license_number, license_authority,
				license_expiry, email, phone, website, country, admin_level_1, admin_level_2, admin_level_3,
				admin_level_4, street, postal_code, latitude, longitude, logo, opening_hours, status
		`, strings.Join(setClauses, ", "), argPos,
	)
	args = append(args, authID)

	row := r.database.QueryRow(ctx, query, args...)

	var pharmacy entities.Pharmacy
	err := row.Scan(
		&pharmacy.PharmacyPublicID, &pharmacy.Name, &pharmacy.LegalName, &pharmacy.Description, &pharmacy.LicenseNumber,
		&pharmacy.LicenseAuthority, &pharmacy.LicenseExpiry, &pharmacy.Email, &pharmacy.Phone, &pharmacy.Website,
		&pharmacy.Country, &pharmacy.AdminLevel1, &pharmacy.AdminLevel2, &pharmacy.AdminLevel3, &pharmacy.AdminLevel4,
		&pharmacy.Street, &pharmacy.PostalCode, &pharmacy.Latitude, &pharmacy.Longitude, &pharmacy.Logo,
		&pharmacy.OpeningHours, &pharmacy.Status,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &pharmacy, nil
}

func (r *pharmacyRepo) UpdateLogo(ctx context.Context, authID int64, logo string) error {
	ctx, span := otel.Tracer(pharmacyErrorTracer).Start(ctx, "UpdateLogo")
	defer span.End()

	query := `
		UPDATE pharmacies
		SET logo = $1, updated_at = NOW()
		WHERE auth_id = $2 AND deleted_at IS NULL
	`

	if err := r.database.Execute(ctx, query, logo, authID); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAuthNotFound, ce.MsgInvalidCredentials, err)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
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
