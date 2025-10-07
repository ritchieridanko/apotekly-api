package validator

import (
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/utils"
)

func nameValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	return len(strings.TrimSpace(value)) >= nameMinLength && len(value) <= nameMaxLength
}

func bioValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	return len(value) <= bioMaxLength
}

func sexValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := utils.Normalize(field.String())
	return value == "male" || value == "female" || value == "other"
}

func birthdateValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value, ok := field.Interface().(time.Time)
	if !ok {
		return false
	}

	return !value.UTC().After(time.Now().UTC())
}

func phoneValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := strings.TrimSpace(field.String())
	return phoneRegex.MatchString(value)
}

func labelValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	return len(strings.TrimSpace(value)) >= labelMinLength && len(value) <= labelMaxLength
}

func notesValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	return len(value) <= notesMaxLength
}

func countryValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := utils.ToTitlecase(field.String())
	return slices.Contains(countries, value)
}

func subdivisionValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	return len(value) <= subdivisionMaxLength
}

func streetValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.String()
	return len(strings.TrimSpace(value)) >= streetMinLength && len(value) <= streetMaxLength
}

func postalCodeValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := strings.TrimSpace(field.String())
	return postalCodeRegex.MatchString(value)
}

func latitudeValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.Float()
	return value >= minLatitude && value <= maxLatitude
}

func longitudeValidator(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := field.Float()
	return value >= minLongitude && value <= maxLongitude
}
