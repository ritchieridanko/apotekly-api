package utils

import "fmt"

func GenerateDynamicRedisKey(prefix string, value any) string {
	return fmt.Sprintf("%s:%v", prefix, value)
}
