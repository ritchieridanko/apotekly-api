package utils

import "fmt"

func GenerateDynamicRedisKey(prefix string, value any) (key string) {
	return fmt.Sprintf("%s:%v", prefix, value)
}
