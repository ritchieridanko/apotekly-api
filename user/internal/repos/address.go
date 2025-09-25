package repos

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"go.opentelemetry.io/otel"
)

const addressErrorTracer string = "repo.address"

type AddressRepo interface {
	Create(ctx context.Context, authID int64, data *entities.NewAddress) (address *entities.Address, err error)
	GetAll(ctx context.Context, authID int64) (addresses []entities.Address, err error)
	Delete(ctx context.Context, authID, addressID int64) (deletedID int64, err error)
	HasPrimary(ctx context.Context, authID int64) (exists bool, err error)
	SetLastAsPrimary(ctx context.Context, authID int64) (newPrimaryID int64, err error)
	UnsetPrimary(ctx context.Context, authID int64) (err error)
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

func (r *addressRepo) Delete(ctx context.Context, authID, addressID int64) (int64, error) {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "Delete")
	defer span.End()

	query := "DELETE FROM addresses WHERE id = $1 AND auth_id = $2 RETURNING id"

	row := r.database.QueryRow(ctx, query, addressID, authID)

	var deletedID int64
	if err := row.Scan(&deletedID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, ce.NewError(span, ce.CodeAddressNotFound, "Address does not exist.", err)
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

	var newPrimaryID int64
	if err := row.Scan(&newPrimaryID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		return 0, ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return newPrimaryID, nil
}

func (r *addressRepo) UnsetPrimary(ctx context.Context, authID int64) error {
	ctx, span := otel.Tracer(addressErrorTracer).Start(ctx, "UnsetPrimary")
	defer span.End()

	query := `
		UPDATE addresses
		SET is_primary = FALSE, updated_at = NOW()
		WHERE auth_id = $1 AND is_primary = TRUE
	`

	if err := r.database.Execute(ctx, query, authID); err != nil {
		return ce.NewError(span, ce.CodeDBQueryExecution, ce.MsgInternalServer, err)
	}

	return nil
}
