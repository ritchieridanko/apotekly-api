package validator

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("email", emailValidator); err != nil {
			return fmt.Errorf("failed to register email validator: %w", err)
		}
		if err := v.RegisterValidation("password", passwordValidator); err != nil {
			return fmt.Errorf("failed to register password validator: %w", err)
		}
		return nil
	}
	return fmt.Errorf("failed to register validators")
}
