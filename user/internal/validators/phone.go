package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	phoneRegex = regexp.MustCompile(`^\d{9,15}$`)
)

func PhoneValidator(fl validator.FieldLevel) bool {
	phone := fl.Field()
	if phone.IsNil() {
		return true
	}
	value, ok := phone.Interface().(*string)
	if !ok {
		return false
	}
	return phoneRegex.MatchString(*value)
}
