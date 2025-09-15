package validators

import "github.com/go-playground/validator/v10"

const (
	maxLength int = 100
)

func BioValidator(fl validator.FieldLevel) bool {
	bio := fl.Field()
	if bio.IsNil() {
		return true
	}
	value, ok := bio.Interface().(*string)
	if !ok {
		return false
	}
	return len(*value) <= maxLength
}
