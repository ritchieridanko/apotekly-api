package dto

import "time"

type CreateAddressRequest struct {
	Recipient    string  `json:"recipient" binding:"required,name"`
	Phone        string  `json:"phone" binding:"required,phone"`
	Label        string  `json:"label" binding:"required,label"`
	Notes        *string `json:"notes" binding:"omitempty,notes"`
	IsPrimary    bool    `json:"is_primary"`
	Country      string  `json:"country" binding:"required,country"`
	Subdivision1 *string `json:"subdivision_1" binding:"omitempty,subdivision"`
	Subdivision2 *string `json:"subdivision_2" binding:"omitempty,subdivision"`
	Subdivision3 *string `json:"subdivision_3" binding:"omitempty,subdivision"`
	Subdivision4 *string `json:"subdivision_4" binding:"omitempty,subdivision"`
	Street       string  `json:"street" binding:"required,street"`
	PostalCode   string  `json:"postal_code" binding:"required,postal_code"`
	Latitude     float64 `json:"latitude" binding:"required,latitude"`
	Longitude    float64 `json:"longitude" binding:"required,longitude"`
}

type UpdateAddressRequest struct {
	Recipient    *string  `json:"recipient" binding:"omitempty,name"`
	Phone        *string  `json:"phone" binding:"omitempty,phone"`
	Label        *string  `json:"label" binding:"omitempty,label"`
	Notes        *string  `json:"notes" binding:"omitempty,notes"`
	Country      *string  `json:"country" binding:"omitempty,country"`
	Subdivision1 *string  `json:"subdivision_1" binding:"omitempty,subdivision"`
	Subdivision2 *string  `json:"subdivision_2" binding:"omitempty,subdivision"`
	Subdivision3 *string  `json:"subdivision_3" binding:"omitempty,subdivision"`
	Subdivision4 *string  `json:"subdivision_4" binding:"omitempty,subdivision"`
	Street       *string  `json:"street" binding:"omitempty,street"`
	PostalCode   *string  `json:"postal_code" binding:"omitempty,postal_code"`
	Latitude     *float64 `json:"latitude" binding:"omitempty,latitude"`
	Longitude    *float64 `json:"longitude" binding:"omitempty,longitude"`
}

type AddressResponse struct {
	ID           int64     `json:"id"`
	Recipient    string    `json:"recipient"`
	Phone        string    `json:"phone"`
	Label        string    `json:"label"`
	Notes        *string   `json:"notes"`
	IsPrimary    bool      `json:"is_primary"`
	Country      string    `json:"country"`
	Subdivision1 *string   `json:"subdivision_1"`
	Subdivision2 *string   `json:"subdivision_2"`
	Subdivision3 *string   `json:"subdivision_3"`
	Subdivision4 *string   `json:"subdivision_4"`
	Street       string    `json:"street"`
	PostalCode   string    `json:"postal_code"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateAddressResponse struct {
	Created AddressResponse  `json:"created"`
	Updated *AddressResponse `json:"updated,omitempty"`
}

type UpdateAddressResponse struct {
	Updated AddressResponse `json:"updated"`
}

type SetPrimaryAddressResponse struct {
	NewPrimaryAddress AddressResponse `json:"new_primary_address"`
	OldPrimaryAddress AddressResponse `json:"old_primary_address"`
}

type DeleteAddressResponse struct {
	NewPrimaryAddress AddressResponse `json:"new_primary_address"`
}
