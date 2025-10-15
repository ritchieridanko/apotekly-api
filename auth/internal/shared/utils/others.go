package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func Normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func ToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func ToInt64Any(value any) (int64, error) {
	v, ok := value.(string)
	if !ok {
		return 0, errors.New("unable to convert value to int64")
	}
	return ToInt64(v)
}
