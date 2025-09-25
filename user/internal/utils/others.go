package utils

import (
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func ToTitlecase(name string) string {
	names := strings.Fields(name)
	if len(names) == 0 {
		return ""
	}

	switch strings.ToLower(names[0]) {
	case "dki":
		return "DKI " + titleCaser.String(strings.Join(names[1:], " "))
	case "di":
		return "DI " + titleCaser.String(strings.Join(names[1:], " "))
	default:
		return titleCaser.String(strings.Join(names, " "))
	}
}

func ToInt64(value string) (number int64, err error) {
	return strconv.ParseInt(value, 10, 64)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
