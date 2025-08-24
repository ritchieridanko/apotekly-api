package utils

import "strings"

func NormalizeString(value string) (normalizedString string) {
	return strings.ToLower(strings.TrimSpace(value))
}
