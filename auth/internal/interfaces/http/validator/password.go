package validator

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	passwordMinLength int = 8
	passwordMaxLength int = 50
)

var (
	lowercaseRegex        = regexp.MustCompile(`[a-z]`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]`)
	numberRegex           = regexp.MustCompile(`[0-9]`)
	specialCharacterRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

func passwordValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	if len(value) < passwordMinLength || len(value) > passwordMaxLength {
		return false
	}

	return lowercaseRegex.MatchString(value) &&
		uppercaseRegex.MatchString(value) &&
		numberRegex.MatchString(value) &&
		specialCharacterRegex.MatchString(value)
}
