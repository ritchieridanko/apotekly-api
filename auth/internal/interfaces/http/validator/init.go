package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("email", emailValidator)
	v.RegisterValidation("password", passwordValidator)

	return &Validator{validator: v}
}

func (v *Validator) Validate(value interface{}) (string, error) {
	err := v.validator.Struct(value)
	if err == nil {
		return "", nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		fe := errs[0]
		return v.toExternal(fe), v.toInternal(fe)
	}

	return ce.MsgInternalServer, err
}

func (v *Validator) toInternal(fe validator.FieldError) error {
	return fmt.Errorf("field '%s' failed on tag '%s'", fe.Field(), fe.Tag())
}

func (v *Validator) toExternal(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s is not a valid email", fe.Value())
	case "password":
		return fmt.Sprintf("%s is not a valid password", fe.Value())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
