package validator

import "regexp"

const (
	nameMinLength        int = 3
	nameMaxLength        int = 100
	bioMaxLength         int = 200
	labelMinLength       int = 3
	labelMaxLength       int = 20
	notesMaxLength       int = 100
	subdivisionMaxLength int = 50
	streetMinLength      int = 5
	streetMaxLength      int = 100

	minLatitude  float64 = -90
	maxLatitude  float64 = 90
	minLongitude float64 = -180
	maxLongitude float64 = 180
)

var (
	phoneRegex      = regexp.MustCompile(`^(?:62|0)8\d{7,11}$`)
	postalCodeRegex = regexp.MustCompile(`^[1-9][0-9]{4}$`)

	countries = []string{
		"Indonesia",
		"United States",
	}
)
