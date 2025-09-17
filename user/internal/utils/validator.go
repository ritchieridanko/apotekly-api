package utils

import (
	"regexp"
	"strings"
	"time"

	"github.com/ritchieridanko/apotekly-api/user/internal/dto"
)

func ValidateNewUser(request dto.ReqNewUser) (isValid bool) {
	if len(strings.TrimSpace(request.Name)) < 3 || len(request.Name) > 50 {
		return false
	}
	if request.Bio != nil && len(*request.Bio) > 100 {
		return false
	}
	if request.Sex != nil && (*request.Sex < 0 || *request.Sex > 2) {
		return false
	}
	if request.Birthdate != nil && request.Birthdate.UTC().After(time.Now().UTC()) {
		return false
	}
	if request.Phone != nil && !regexp.MustCompile(`^\d{9,15}$`).MatchString(*request.Phone) {
		return false
	}
	return true
}
