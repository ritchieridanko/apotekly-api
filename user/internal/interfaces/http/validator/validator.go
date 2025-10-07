package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("name", nameValidator)
	v.RegisterValidation("bio", bioValidator)
	v.RegisterValidation("sex", sexValidator)
	v.RegisterValidation("birthdate", birthdateValidator)
	v.RegisterValidation("phone", phoneValidator)
	v.RegisterValidation("label", labelValidator)
	v.RegisterValidation("notes", notesValidator)
	v.RegisterValidation("country", countryValidator)
	v.RegisterValidation("subdivision", subdivisionValidator)
	v.RegisterValidation("street", streetValidator)
	v.RegisterValidation("postal_code", postalCodeValidator)
	v.RegisterValidation("latitude", latitudeValidator)
	v.RegisterValidation("longitude", longitudeValidator)

	return &Validator{validator: v}
}

func (v *Validator) Validate(value interface{}) (string, error) {
	err := v.validator.Struct(value)
	if err == nil {
		return "", nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		fe := errs[0]
		return v.toExternal(fe), fmt.Errorf("failed to validate: %w", v.toInternal(fe))
	}

	return ce.MsgInternalServer, fmt.Errorf("failed to validate: %w", err)
}

func (v *Validator) toInternal(fe validator.FieldError) error {
	return fmt.Errorf("field '%s' failed on tag '%s'", fe.Field(), fe.Tag())
}

func (v *Validator) toExternal(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "name":
		return fmt.Sprintf("Name must be between %d and %d characters", nameMinLength, nameMaxLength)
	case "bio":
		return fmt.Sprintf("Bio must not exceed %d characters", bioMaxLength)
	case "sex":
		return fmt.Sprintf("%s is not a valid sex option", fe.Value())
	case "birthdate":
		return fmt.Sprintf("%s is not a valid birthdate", fe.Value())
	case "phone":
		return fmt.Sprintf("%s is not a valid phone number", fe.Value())
	case "label":
		return fmt.Sprintf("Label must be between %d and %d characters", labelMinLength, labelMaxLength)
	case "notes":
		return fmt.Sprintf("Notes must not exceed %d characters", notesMaxLength)
	case "country":
		return fmt.Sprintf("%s is not a valid country option", fe.Value())
	case "subdivision":
		return fmt.Sprintf("Subdivision must not exceed %d characters", subdivisionMaxLength)
	case "street":
		return fmt.Sprintf("Street must be between %d and %d characters", streetMinLength, streetMaxLength)
	case "postal_code":
		return fmt.Sprintf("%s is not a valid postal code", fe.Value())
	case "latitude":
		return fmt.Sprintf("%s is not a valid latitude", fe.Value())
	case "longitude":
		return fmt.Sprintf("%s is not a valid longitude", fe.Value())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
