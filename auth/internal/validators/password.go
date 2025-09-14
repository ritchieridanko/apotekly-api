package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	lowercaseRegex        = regexp.MustCompile(`[a-z]`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]`)
	numberRegex           = regexp.MustCompile(`[0-9]`)
	specialCharacterRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

const (
	minPasswordLength int = 8
	maxPasswordLength int = 50
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return false
	}

	return lowercaseRegex.MatchString(password) && uppercaseRegex.MatchString(password) && numberRegex.MatchString(password) && specialCharacterRegex.MatchString(password)
}
