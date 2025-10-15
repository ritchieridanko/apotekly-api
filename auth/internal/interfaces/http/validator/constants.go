package validator

import "regexp"

const (
	passwordMinLength int = 8
	passwordMaxLength int = 50
)

var (
	emailRegex            = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	lowercaseRegex        = regexp.MustCompile(`[a-z]`)
	uppercaseRegex        = regexp.MustCompile(`[A-Z]`)
	numberRegex           = regexp.MustCompile(`[0-9]`)
	specialCharacterRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)
