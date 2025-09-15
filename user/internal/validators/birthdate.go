package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func BirthdateValidator(fl validator.FieldLevel) bool {
	birthdate := fl.Field()
	if birthdate.IsNil() {
		return true
	}
	_, ok := birthdate.Interface().(*time.Time)
	return ok
}
