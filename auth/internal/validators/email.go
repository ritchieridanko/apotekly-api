package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func emailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return emailRegex.MatchString(email)
}
