package repos

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"go.opentelemetry.io/otel"
)

const addressErrorTracer string = "repo.address"

type AddressRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewAddress) (address *entities.Address, err error)
	GetAll(ctx context.Context, authID int64) (addresses []entities.Address, err error)
	Update(ctx context.Context, authID, addressID int64, data *entities.AddressChange) (address *entities.Address, err error)
	Delete(ctx context.Context, authID, addressID int64) (deletedID int64, err error)
	HasPrimary(ctx context.Context, authID int64) (exists bool, err error)
	SetAsPrimary(ctx context.Context, authID, addressID int64) (newPrimaryID int64, err error)
	SetLastAsPrimary(ctx context.Context, authID int64) (newPrimaryID int64, err error)
	UnsetPrimary(ctx context.Context, authID int64) (unsetPrimaryID int64, err error)
}

type addressRepo struct {
	database db.DBService
}

func NewAddressRepo(database db.DBService) AddressRepo {
	return &addressRepo{database}
}

func (r *addressRepo) Create(ctx context.Context, authID int64, data *entities.NewAddress) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Create")
	defer span.End()

	query := `
		INSERT INTO
			addresses (
				auth_id, receiver, phone, label, notes, is_primary, country,
				admin_level_1, admin_level_2, admin_level_3, admin_level_4,
				street, postal_code, latitude, longitude, location
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, ST_SetSRID(ST_MakePoint($15, $14), 4326)
		)
		RETURNING
			id, receiver, phone, label, notes, is_primary, country,
			admin_level_1, admin_level_2, admin_level_3, admin_level_4,
			street, postal_code, latitude, longitude
	`

	row := r.database.QueryRow(
		ctx, query, authID,
		data.Receiver, data.Phone, data.Label, data.Notes, data.IsPrimary, data.Country,
		data.AdminLevel1, data.AdminLevel2, data.AdminLevel3, data.AdminLevel4,
		data.Street, data.PostalCode, data.Latitude, data.Longitude,
	)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Receiver, &address.Phone, &address.Label, &address.Notes,
		&address.IsPrimary, &address.Country, &address.AdminLevel1, &address.AdminLevel2,
		&address.AdminLevel3, &address.AdminLevel4, &address.Street, &address.PostalCode,
		&address.Latitude, &address.Longitude,
	)
	if err != nil {
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &address, nil
}

func (r *addressRepo) GetAll(ctx context.Context, authID int64) ([]entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "GetAll")
	defer span.End()

	query := `
		SELECT
			id, receiver, phone, label, notes, is_primary, country,
			admin_level_1, admin_level_2, admin_level_3, admin_level_4,
			street, postal_code, latitude, longitude
		FROM
			addresses
		WHERE
			auth_id = $1
		ORDER BY
			is_primary DESC,
			created_at DESC
	`

	rows, err := r.database.QueryAll(ctx, query, authID)
	if err != nil {
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}
	defer rows.Close()

	addresses := make([]entities.Address, 0)
	for rows.Next() {
		var address entities.Address

		err := rows.Scan(
			&address.ID, &address.Receiver, &address.Phone, &address.Label, &address.Notes,
			&address.IsPrimary, &address.Country, &address.AdminLevel1, &address.AdminLevel2,
			&address.AdminLevel3, &address.AdminLevel4, &address.Street, &address.PostalCode,
			&address.Latitude, &address.Longitude,
		)
		if err != nil {
			return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
		}

		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	if len(addresses) == 0 {
		return []entities.Address{}, nil
	}

	return addresses, nil
}

func (r *addressRepo) Update(ctx context.Context, authID, addressID int64, data *entities.AddressChange) (*entities.Address, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Update")
	defer span.End()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if data.Receiver != nil {
		setClauses = append(setClauses, fmt.Sprintf("receiver = $%d", argPos))
		args = append(args, *data.Receiver)
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
	if data.IsPrimary != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_primary = $%d", argPos))
		args = append(args, *data.IsPrimary)
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
	if len(setClauses) == 0 {
		return nil, ce.NewError(span, ce.CodeInvalidPayload, ce.MsgNoFieldsToUpdate, ce.ErrNoFieldsProvided)
	}
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
			UPDATE addresses
			SET %s
			WHERE id = $%d AND auth_id = $%d
			RETURNING
				id, receiver, phone, label, notes, is_primary, country,
				admin_level_1, admin_level_2, admin_level_3, admin_level_4,
				street, postal_code, latitude, longitude
		`, strings.Join(setClauses, ", "), argPos, argPos+1,
	)
	args = append(args, addressID, authID)

	row := r.database.QueryRow(ctx, query, args...)

	var address entities.Address
	err := row.Scan(
		&address.ID, &address.Receiver, &address.Phone, &address.Label, &address.Notes,
		&address.IsPrimary, &address.Country, &address.AdminLevel1, &address.AdminLevel2,
		&address.AdminLevel3, &address.AdminLevel4, &address.Street, &address.PostalCode,
		&address.Latitude, &address.Longitude,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, err)
		}
		return nil, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return &address, nil
}

func (r *addressRepo) Delete(ctx context.Context, authID, addressID int64) (int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Delete")
	defer span.End()

	query := "DELETE FROM addresses WHERE id = $1 AND auth_id = $2 RETURNING id"

	row := r.database.QueryRow(ctx, query, addressID, authID)

	var deletedID int64
	if err := row.Scan(&deletedID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, err)
		}
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return deletedID, nil
}

func (r *addressRepo) HasPrimary(ctx context.Context, authID int64) (bool, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "HasPrimary")
	defer span.End()

	query := `
		SELECT 1
		FROM addresses
		WHERE auth_id = $1 AND is_primary = TRUE
	`
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

func (r *addressRepo) SetAsPrimary(ctx context.Context, authID, addressID int64) (int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "SetAsPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = TRUE, updated_at = NOW()
		WHERE id = $1 AND auth_id = $2
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, addressID, authID)

	var newPrimaryID int64
	if err := row.Scan(&newPrimaryID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, ce.NewError(span, ce.CodeAddressNotFound, ce.MsgAddressNotFound, err)
		}
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return newPrimaryID, nil
}

func (r *addressRepo) SetLastAsPrimary(ctx context.Context, authID int64) (int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "SetLastAsPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = TRUE, updated_at = NOW()
		WHERE id = (
			SELECT id
			FROM addresses
			WHERE auth_id = $1
			ORDER BY created_at DESC
			LIMIT 1
		)
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, authID)

	var addressID int64
	if err := row.Scan(&addressID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return addressID, nil
}

func (r *addressRepo) UnsetPrimary(ctx context.Context, authID int64) (int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "UnsetPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = FALSE, updated_at = NOW()
		WHERE auth_id = $1 AND is_primary = TRUE
		RETURNING id
	`

	row := r.database.QueryRow(ctx, query, authID)

	var addressID int64
	if err := row.Scan(&addressID); err != nil {
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return addressID, nil
}
