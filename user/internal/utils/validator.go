package utils

import "time"

func ValidatorIsBirthdateValid(birthdate time.Time) (isValid bool) {
	now := time.Now().UTC()
	return birthdate.Before(now)
}

func ValidatorIsSexValid(sex int16) (isValid bool) {
	male := int16(0)
	female := int16(1)
	others := int16(2)
	return sex == male || sex == female || sex == others
}
