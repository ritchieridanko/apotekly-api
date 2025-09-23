package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ritchieridanko/apotekly-api/user/internal/dto"
)

const (
	nameMinLength int = 3
	nameMaxLength int = 50
	bioMaxLength  int = 100
)

var (
	phoneRegex = regexp.MustCompile(`^\d{9,15}$`)

	acceptableFileTypes []string = []string{"image/png", "image/jpeg"}
)

func ValidateNewUser(request dto.ReqNewUser) (errString string) {
	if len(strings.TrimSpace(request.Name)) < nameMinLength || len(request.Name) > nameMaxLength {
		return fmt.Sprintf("Name must be between %d and %d characters.", nameMinLength, nameMaxLength)
	}
	if request.Bio != nil && len(*request.Bio) > bioMaxLength {
		return fmt.Sprintf("Bio must not exceed %d characters.", bioMaxLength)
	}
	if request.Sex != nil && (*request.Sex < 0 || *request.Sex > 2) {
		return "Sex is invalid."
	}
	if request.Birthdate != nil && request.Birthdate.UTC().After(time.Now().UTC()) {
		return "Birthdate is invalid."
	}
	if request.Phone != nil && !phoneRegex.MatchString(*request.Phone) {
		return "Phone is invalid."
	}
	return ""
}

func ValidateUserUpdate(request dto.ReqUserUpdate) (errString string) {
	if request.Name != nil && (len(strings.TrimSpace(*request.Name)) < nameMinLength || len(*request.Name) > nameMaxLength) {
		return fmt.Sprintf("Name must be between %d and %d characters.", nameMinLength, nameMaxLength)
	}
	if request.Bio != nil && len(*request.Bio) > bioMaxLength {
		return fmt.Sprintf("Bio must not exceed %d characters.", bioMaxLength)
	}
	if request.Sex != nil && (*request.Sex < 0 || *request.Sex > 2) {
		return "Sex is invalid."
	}
	if request.Birthdate != nil && request.Birthdate.UTC().After(time.Now().UTC()) {
		return "Birthdate is invalid."
	}
	if request.Phone != nil && !phoneRegex.MatchString(*request.Phone) {
		return "Phone is invalid."
	}
	return ""
}

func ValidateImageFile(imageBuf []byte) (err error) {
	// validate content type from first 512 bytes
	fileType := http.DetectContentType(imageBuf[:min(len(imageBuf), 512)])

	allowed := false
	for _, acceptableType := range acceptableFileTypes {
		if fileType == acceptableType {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("invalid file type: %s", fileType)
	}

	return nil
}
