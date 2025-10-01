package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/constants"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/dto"
)

const (
	nameMaxLength          int = 100
	legalNameMaxLength     int = 200
	descriptionMaxLength   int = 300
	licenseNumberMaxLength int = 100
)

var (
	acceptableFileTypes []string = []string{"image/png", "image/jpeg"}
	days                []string = []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex       = regexp.MustCompile(`^(?:62|0)8\d{7,11}$`)
	openingHourRegex = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d-(?:[01]\d|2[0-3]):[0-5]\d$`)
)

func ValidateNewPharmacy(request dto.ReqNewPharmacy) (errString string) {
	if len(strings.TrimSpace(request.Name)) == 0 {
		return "Pharmacy name must not be empty."
	}
	if len(request.Name) > nameMaxLength {
		return fmt.Sprintf("Pharmacy name must not exceed %d characters.", nameMaxLength)
	}
	if request.LegalName != nil && len(*request.LegalName) > legalNameMaxLength {
		return fmt.Sprintf("Legal name must not exceed %d characters.", legalNameMaxLength)
	}
	if request.Description != nil && len(*request.Description) > descriptionMaxLength {
		return fmt.Sprintf("Description must not exceed %d characters.", descriptionMaxLength)
	}
	if len(strings.TrimSpace(request.LicenseNumber)) == 0 {
		return "License number must not be empty."
	}
	if len(request.LicenseNumber) > licenseNumberMaxLength {
		return fmt.Sprintf("License number must not exceed %d characters.", licenseNumberMaxLength)
	}
	if len(strings.TrimSpace(request.LicenseAuthority)) == 0 {
		return "License authority must not be empty."
	}
	if len(request.LicenseAuthority) > nameMaxLength {
		return fmt.Sprintf("License authority must not exceed %d characters.", nameMaxLength)
	}
	if request.LicenseExpiry != nil && request.LicenseExpiry.UTC().Before(time.Now().UTC()) {
		return "License expiry is invalid."
	}
	if request.Email != nil && !emailRegex.MatchString(*request.Email) {
		return "Email is invalid."
	}
	if request.Phone != nil && !phoneRegex.MatchString(*request.Phone) {
		return "Phone is invalid."
	}
	if request.Website != nil {
		if u, err := url.ParseRequestURI(*request.Website); err != nil || u.Scheme == "" || u.Host == "" {
			return "Website is invalid."
		}
	}
	if validateErr := ValidateOpeningHours(request.OpeningHours); validateErr != "" {
		return validateErr
	}
	return ""
}

func ValidateOpeningHours(data constants.OpeningHours) (errString string) {
	for key, value := range data {
		if exists := slices.Contains(days, key); !exists {
			return fmt.Sprintf("Opening day %s is invalid.", key)
		}

		for _, hour := range value {
			if !openingHourRegex.MatchString(hour) {
				return fmt.Sprintf("Opening hour %s is invalid.", hour)
			}

			var startStr, endStr string
			fmt.Sscanf(hour, "%5s-%5s", &startStr, &endStr)

			start, err := time.Parse("15:04", startStr)
			if err != nil {
				return fmt.Sprintf("Opening hour %s is invalid.", hour)
			}
			end, err := time.Parse("15:04", endStr)
			if err != nil {
				return fmt.Sprintf("Closing hour %s is invalid.", hour)
			}

			if !start.Before(end) {
				return "Start time must be before end time."
			}
		}
	}
	return ""
}

func ValidateImageFile(imageBuf []byte) (err error) {
	// validate content type from first 512 bytes
	fileType := http.DetectContentType(imageBuf[:min(len(imageBuf), 512)])

	allowed := false
	for _, acceptableType := range acceptableFileTypes {
		if fileType == acceptableType {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("invalid file type: %s", fileType)
	}

	return nil
}
