package utils

import "fmt"

func CacheCreateKey(prefix string, value any) (key string) {
	return fmt.Sprintf("%s:%v", prefix, value)
}
