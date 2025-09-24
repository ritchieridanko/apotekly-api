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

	receiverMinLength int = 2
	receiverMaxLength int = 20
	labelMinLength    int = 2
	labelMaxLength    int = 20
	notesMaxLength    int = 50
)

var (
	phoneRegex = regexp.MustCompile(`^(?:62|0)8\d{7,11}$`)

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

func ValidateNewAddress(request dto.ReqNewAddress) (errString string) {
	if len(strings.TrimSpace(request.Receiver)) < receiverMinLength || len(request.Receiver) > receiverMaxLength {
		return fmt.Sprintf("Receiver name must be between %d and %d characters.", receiverMinLength, receiverMaxLength)
	}
	if !phoneRegex.MatchString(request.Phone) {
		return "Phone is invalid."
	}
	if len(strings.TrimSpace(request.Label)) < labelMinLength || len(request.Label) > labelMaxLength {
		return fmt.Sprintf("Label must be between %d and %d characters.", labelMinLength, labelMaxLength)
	}
	if request.Notes != nil && len(*request.Notes) > notesMaxLength {
		return fmt.Sprintf("Notes must not exceed %d characters.", notesMaxLength)
	}
	return ""
}
