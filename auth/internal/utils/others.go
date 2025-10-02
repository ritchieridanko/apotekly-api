package utils

import (
	"strconv"
	"strings"
)

func Normalize(value string) (result string) {
	return strings.ToLower(strings.TrimSpace(value))
}

func ToInt64(value string) (number int64, err error) {
	return strconv.ParseInt(value, 10, 64)
}
