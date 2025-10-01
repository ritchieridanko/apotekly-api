package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/constants"
)

type Pharmacy struct {
	PharmacyPublicID uuid.UUID
	Name             string
	LegalName        *string
	Description      *string
	LicenseNumber    string
	LicenseAuthority string
	LicenseExpiry    *time.Time
	Email            *string
	Phone            *string
	Website          *string
	Country          string
	AdminLevel1      *string
	AdminLevel2      *string
	AdminLevel3      *string
	AdminLevel4      *string
	Street           string
	PostalCode       string
	Latitude         float64
	Longitude        float64
	Logo             *string
	OpeningHours     constants.OpeningHours
	Status           string
}

type NewPharmacy struct {
	PharmacyPublicID uuid.UUID
	Name             string
	LegalName        *string
	Description      *string
	LicenseNumber    string
	LicenseAuthority string
	LicenseExpiry    *time.Time
	Email            *string
	Phone            *string
	Website          *string
	Country          string
	AdminLevel1      *string
	AdminLevel2      *string
	AdminLevel3      *string
	AdminLevel4      *string
	Street           string
	PostalCode       string
	Latitude         float64
	Longitude        float64
	Logo             *string
	OpeningHours     constants.OpeningHours
}
