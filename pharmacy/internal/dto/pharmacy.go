package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/types"
)

type RespPharmacy struct {
	ID               uuid.UUID          `json:"id"`
	Name             string             `json:"name"`
	LegalName        *string            `json:"legal_name"`
	Description      *string            `json:"description"`
	LicenseNumber    string             `json:"license_number"`
	LicenseAuthority string             `json:"license_authority"`
	LicenseExpiry    *time.Time         `json:"license_expiry"`
	Email            *string            `json:"email"`
	Phone            *string            `json:"phone"`
	Website          *string            `json:"website"`
	Country          string             `json:"country"`
	AdminLevel1      *string            `json:"admin_level_1"`
	AdminLevel2      *string            `json:"admin_level_2"`
	AdminLevel3      *string            `json:"admin_level_3"`
	AdminLevel4      *string            `json:"admin_level_4"`
	Street           string             `json:"street"`
	PostalCode       string             `json:"postal_code"`
	Latitude         float64            `json:"latitude"`
	Longitude        float64            `json:"longitude"`
	Logo             *string            `json:"logo"`
	OpeningHours     types.OpeningHours `json:"opening_hours"`
	Status           string             `json:"status"`
}

type ReqNewPharmacy struct {
	Name             string                `form:"name" binding:"required"`
	LegalName        *string               `form:"legal_name"`
	Description      *string               `form:"description"`
	LicenseNumber    string                `form:"license_number" binding:"required"`
	LicenseAuthority string                `form:"license_authority" binding:"required"`
	LicenseExpiry    *time.Time            `form:"license_expiry" time_format:"2006-01-02"`
	Email            *string               `form:"email"`
	Phone            *string               `form:"phone"`
	Website          *string               `form:"website"`
	Country          string                `form:"country" binding:"required"`
	AdminLevel1      *string               `form:"admin_level_1"`
	AdminLevel2      *string               `form:"admin_level_2"`
	AdminLevel3      *string               `form:"admin_level_3"`
	AdminLevel4      *string               `form:"admin_level_4"`
	Street           string                `form:"street" binding:"required"`
	PostalCode       string                `form:"postal_code" binding:"required"`
	Latitude         float64               `form:"latitude" binding:"required"`
	Longitude        float64               `form:"longitude" binding:"required"`
	OpeningHours     types.OpeningHours    `form:"opening_hours" binding:"required"`
	Image            *multipart.FileHeader `form:"image"`
}

type RespNewPharmacy struct {
	Created RespPharmacy `json:"created"`
}

type ReqUpdatePharmacy struct {
	Name             *string             `json:"name"`
	LegalName        *string             `json:"legal_name"`
	Description      *string             `json:"description"`
	LicenseNumber    *string             `json:"license_number"`
	LicenseAuthority *string             `json:"license_authority"`
	LicenseExpiry    *time.Time          `json:"license_expiry"`
	Email            *string             `json:"email"`
	Phone            *string             `json:"phone"`
	Website          *string             `json:"website"`
	Country          *string             `json:"country"`
	AdminLevel1      *string             `json:"admin_level_1"`
	AdminLevel2      *string             `json:"admin_level_2"`
	AdminLevel3      *string             `json:"admin_level_3"`
	AdminLevel4      *string             `json:"admin_level_4"`
	Street           *string             `json:"street"`
	PostalCode       *string             `json:"postal_code"`
	Latitude         *float64            `json:"latitude"`
	Longitude        *float64            `json:"longitude"`
	OpeningHours     *types.OpeningHours `json:"opening_hours"`
}

type RespUpdatePharmacy struct {
	Updated RespPharmacy `json:"updated"`
}
