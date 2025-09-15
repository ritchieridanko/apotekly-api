package validators

import "github.com/go-playground/validator/v10"

func SexValidator(fl validator.FieldLevel) bool {
	sex := fl.Field()
	if sex.IsNil() {
		return true
	}
	value, ok := sex.Interface().(*int16)
	if !ok {
		return false
	}
	return *value >= 0
}
