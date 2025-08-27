package validators

import "github.com/go-playground/validator/v10"

func Initialize(validator *validator.Validate) error {
	if err := validator.RegisterValidation("bio", BioValidator); err != nil {
		return err
	}
	if err := validator.RegisterValidation("sex", SexValidator); err != nil {
		return err
	}
	if err := validator.RegisterValidation("birthdate", BirthdateValidator); err != nil {
		return err
	}
	if err := validator.RegisterValidation("phone", PhoneValidator); err != nil {
		return err
	}
	return nil
}
