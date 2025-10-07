package utils

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func NormalizePtr(value *string) *string {
	if value == nil {
		return nil
	}
	result := Normalize(*value)
	return &result
}

func ToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func ToTitlecase(value string) string {
	values := strings.Fields(value)
	if len(values) == 0 {
		return ""
	}

	switch strings.ToLower(values[0]) {
	case "dki", "di":
		return strings.ToUpper(values[0]) + " " + titleCaser.String(strings.Join(values[1:], " "))
	default:
		return titleCaser.String(strings.Join(values, " "))
	}
}

func ToTitlecasePtr(value *string) *string {
	if value == nil {
		return nil
	}
	result := ToTitlecase(*value)
	return &result
}

func TrimSpacePtr(value *string) *string {
	if value == nil {
		return nil
	}
	result := strings.TrimSpace(*value)
	return &result
}
