package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/database"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"go.opentelemetry.io/otel"
)

const addressErrorTracer string = "repository.address"

type AddressRepository interface {
	Create(ctx context.Context, authID int64, data *entities.CreateAddress) (address *entities.Address, err error)
	GetAll(ctx context.Context, authID int64) (addresses []entities.Address, err error)
	Update(ctx context.Context, authID, addressID int64, data *entities.UpdateAddress) (address *entities.Address, err error)
	Delete(ctx context.Context, authID, addressID int64) (err error)
	HasPrimary(ctx context.Context, authID int64) (exists bool, err error)
	SetPrimary(ctx context.Context, authID, addressID int64) (address *entities.Address, err error)
	UnsetPrimary(ctx context.Context, authID int64) (address *entities.Address, err error)
	SetLastUpdatedPrimary(ctx context.Context, authID int64) (address *entities.Address, err error)
}

type addressRepository struct {
	database *database.Database
}

func NewAddressRepository(database *database.Database) AddressRepository {
	return &addressRepository{database}
}

func (r *addressRepository) Create(ctx context.Context, authID int64, data *entities.CreateAddress) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO addresses (
			auth_id, recipient, phone, label, notes, is_primary, country,
			subdivision_1, subdivision_2, subdivision_3, subdivision_4,
			street, postal_code, latitude, longitude, location
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, ST_SetSRID(ST_MakePoint($15, $14), 4326)
		)
		RETURNING
			address_id, recipient, phone, label, notes, is_primary,
			country, subdivision_1, subdivision_2, subdivision_3,
			subdivision_4, street, postal_code, latitude, longitude,
			created_at, updated_at
	`

	row := r.database.QueryRow(
		ctx, query,
		authID, data.Recipient, data.Phone, data.Label, data.Notes, data.IsPrimary,
		data.Country, data.Subdivision1, data.Subdivision2, data.Subdivision3,
		data.Subdivision4, data.Street, data.PostalCode, data.Latitude, data.Longitude,
	)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label,
		&address.Notes, &address.IsPrimary, &address.Country, &address.Subdivision1,
		&address.Subdivision2, &address.Subdivision3, &address.Subdivision4,
		&address.Street, &address.PostalCode, &address.Latitude, &address.Longitude,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create address: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &address, nil
}

func (r *addressRepository) GetAll(ctx context.Context, authID int64) ([]entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "GetAll")
	defer span.End()

	query := `
		SELECT
			address_id, recipient, phone, label, notes, is_primary,
			country, subdivision_1, subdivision_2, subdivision_3,
			subdivision_4, street, postal_code, latitude, longitude,
			created_at, updated_at
		FROM
			addresses
		WHERE
			auth_id = $1
		ORDER BY
			is_primary DESC,
			updated_at DESC
	`

	rows, err := r.database.QueryAll(ctx, query, authID)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch all addresses: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}
	defer rows.Close()

	addresses := make([]entities.Address, 0)
	for rows.Next() {
		var address entities.Address

		err := rows.Scan(
			&address.ID, &address.Recipient, &address.Phone, &address.Label,
			&address.Notes, &address.IsPrimary, &address.Country, &address.Subdivision1,
			&address.Subdivision2, &address.Subdivision3, &address.Subdivision4,
			&address.Street, &address.PostalCode, &address.Latitude, &address.Longitude,
			&address.CreatedAt, &address.UpdatedAt,
		)
		if err != nil {
			wErr := fmt.Errorf("failed to fetch all addresses: %w", err)
			return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		wErr := fmt.Errorf("failed to fetch all addresses: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	if len(addresses) == 0 {
		return []entities.Address{}, nil
	}

	return addresses, nil
}

func (r *addressRepository) Update(ctx context.Context, authID, addressID int64, data *entities.UpdateAddress) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Update")
	defer span.End()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if data.Recipient != nil {
		setClauses = append(setClauses, fmt.Sprintf("recipient = $%d", argPos))
		args = append(args, *data.Recipient)
		argPos++
	}
	if data.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *data.Phone)
		argPos++
	}
	if data.Label != nil {
		setClauses = append(setClauses, fmt.Sprintf("label = $%d", argPos))
		args = append(args, *data.Label)
		argPos++
	}
	if data.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argPos))
		args = append(args, *data.Notes)
		argPos++
	}
	if data.Country != nil {
		setClauses = append(setClauses, fmt.Sprintf("country = $%d", argPos))
		args = append(args, *data.Country)
		argPos++
	}
	if data.Subdivision1 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_1 = $%d", argPos))
		args = append(args, *data.Subdivision1)
		argPos++
	}
	if data.Subdivision2 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_2 = $%d", argPos))
		args = append(args, *data.Subdivision2)
		argPos++
	}
	if data.Subdivision3 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_3 = $%d", argPos))
		args = append(args, *data.Subdivision3)
		argPos++
	}
	if data.Subdivision4 != nil {
		setClauses = append(setClauses, fmt.Sprintf("subdivision_4 = $%d", argPos))
		args = append(args, *data.Subdivision4)
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
	if len(setClauses) == 0 {
		err := fmt.Errorf("failed to update address: %w", ce.ErrNoFieldsProvided)
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgNoFieldsToUpdate, err)
	}
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
			UPDATE addresses
			SET %s
			WHERE address_id = $%d AND auth_id = $%d
			RETURNING
				address_id, recipient, phone, label, notes, is_primary,
				country, subdivision_1, subdivision_2, subdivision_3,
				subdivision_4, street, postal_code, latitude, longitude,
				created_at, updated_at
		`, strings.Join(setClauses, ", "), argPos, argPos+1,
	)
	args = append(args, addressID, authID)

	row := r.database.QueryRow(ctx, query, args...)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label,
		&address.Notes, &address.IsPrimary, &address.Country, &address.Subdivision1,
		&address.Subdivision2, &address.Subdivision3, &address.Subdivision4,
		&address.Street, &address.PostalCode, &address.Latitude, &address.Longitude,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to update address: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &address, nil
}

func (r *addressRepository) Delete(ctx context.Context, authID, addressID int64) error {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Delete")
	defer span.End()

	query := "DELETE FROM addresses WHERE address_id = $1 AND auth_id = $2"

	if err := r.database.Execute(ctx, query, addressID, authID); err != nil {
		wErr := fmt.Errorf("failed to delete address: %w", err)
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, wErr)
		}
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return nil
}

func (r *addressRepository) HasPrimary(ctx context.Context, authID int64) (bool, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "HasPrimary")
	defer span.End()

	query := "SELECT 1 FROM addresses WHERE auth_id = $1 AND is_primary = TRUE"
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, authID)

	var exists int
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		wErr := fmt.Errorf("failed to fetch address: %w", err)
		return false, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return true, nil
}

func (r *addressRepository) SetPrimary(ctx context.Context, authID, addressID int64) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "SetPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = TRUE, updated_at = NOW()
		WHERE address_id = $1 AND auth_id = $2
		RETURNING
			address_id, recipient, phone, label, notes, is_primary,
			country, subdivision_1, subdivision_2, subdivision_3,
			subdivision_4, street, postal_code, latitude, longitude,
			created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, addressID, authID)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label,
		&address.Notes, &address.IsPrimary, &address.Country, &address.Subdivision1,
		&address.Subdivision2, &address.Subdivision3, &address.Subdivision4,
		&address.Street, &address.PostalCode, &address.Latitude, &address.Longitude,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to set primary: %w", err)
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, wErr)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &address, nil
}

func (r *addressRepository) UnsetPrimary(ctx context.Context, authID int64) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "UnsetPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = FALSE, updated_at = NOW()
		WHERE auth_id = $1 AND is_primary = TRUE
		RETURNING
			address_id, recipient, phone, label, notes, is_primary,
			country, subdivision_1, subdivision_2, subdivision_3,
			subdivision_4, street, postal_code, latitude, longitude,
			created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, authID)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label,
		&address.Notes, &address.IsPrimary, &address.Country, &address.Subdivision1,
		&address.Subdivision2, &address.Subdivision3, &address.Subdivision4,
		&address.Street, &address.PostalCode, &address.Latitude, &address.Longitude,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		wErr := fmt.Errorf("failed to unset primary: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &address, nil
}

func (r *addressRepository) SetLastUpdatedPrimary(ctx context.Context, authID int64) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "SetLastUpdatedPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = TRUE, updated_at = NOW()
		WHERE address_id = (
			SELECT address_id FROM addresses WHERE auth_id = $1
			ORDER BY updated_at DESC LIMIT 1
		)
		RETURNING
			address_id, recipient, phone, label, notes, is_primary,
			country, subdivision_1, subdivision_2, subdivision_3,
			subdivision_4, street, postal_code, latitude, longitude,
			created_at, updated_at
	`

	row := r.database.QueryRow(ctx, query, authID)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Recipient, &address.Phone, &address.Label,
		&address.Notes, &address.IsPrimary, &address.Country, &address.Subdivision1,
		&address.Subdivision2, &address.Subdivision3, &address.Subdivision4,
		&address.Street, &address.PostalCode, &address.Latitude, &address.Longitude,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, nil
		}
		wErr := fmt.Errorf("failed to set last updated primary: %w", err)
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, wErr)
	}

	return &address, nil
}
