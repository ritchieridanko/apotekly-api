package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func emailValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := strings.TrimSpace(field.String())
	return emailRegex.MatchString(value)
}

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

	return lowercaseRegex.MatchString(value) && uppercaseRegex.MatchString(value) && numberRegex.MatchString(value) && specialCharacterRegex.MatchString(value)
}
