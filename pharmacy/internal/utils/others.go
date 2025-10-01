package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func Normalize(value string) (result string) {
	return strings.ToLower(strings.TrimSpace(value))
}

func NormalizePtr(value *string) (result *string) {
	if value == nil {
		return nil
	}
	normalizedValue := Normalize(*value)
	return &normalizedValue
}

func ToTitlecase(value string) (result string) {
	values := strings.Fields(value)
	if len(values) == 0 {
		return ""
	}

	switch strings.ToLower(values[0]) {
	case "dki":
		return "DKI " + titleCaser.String(strings.Join(values[1:], " "))
	case "di":
		return "DI " + titleCaser.String(strings.Join(values[1:], " "))
	default:
		return titleCaser.String(strings.Join(values, " "))
	}
}

func ToTitlecasePtr(value *string) (result *string) {
	if value == nil {
		return nil
	}
	titlecasedValue := ToTitlecase(*value)
	return &titlecasedValue
}

func TrimSpacePtr(value *string) (result *string) {
	if value == nil {
		return nil
	}
	trimmedValue := strings.TrimSpace(*value)
	return &trimmedValue
}
