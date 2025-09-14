package validators

import "github.com/go-playground/validator/v10"

func Initialize(validator *validator.Validate) error {
	if err := validator.RegisterValidation("email", EmailValidator); err != nil {
		return err
	}
	if err := validator.RegisterValidation("password", PasswordValidator); err != nil {
		return err
	}
	return nil
}
