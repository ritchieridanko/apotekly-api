package config

type cacheConfig struct {
	Host       string
	Port       string
	Pass       string
	MaxRetries int
	BaseDelay  int
}

var cacheCfg *cacheConfig

func loadCacheConfig() {
	cacheCfg = &cacheConfig{
		Host:       getEnv("CACHE_HOST"),
		Port:       getEnv("CACHE_PORT"),
		Pass:       getEnv("CACHE_PASS"),
		MaxRetries: getNumberEnvWithFallback("CACHE_MAX_RETRIES", 3),
		BaseDelay:  getNumberEnvWithFallback("CACHE_BASE_DELAY", 100),
	}
}

func CacheGetHost() (host string) {
	return cacheCfg.Host
}

func CacheGetPort() (port string) {
	return cacheCfg.Port
}

func CacheGetPass() (pass string) {
	return cacheCfg.Pass
}

func CacheGetMaxRetries() (max int) {
	return cacheCfg.MaxRetries
}

func CacheGetBaseDelay() (delay int) {
	return cacheCfg.BaseDelay
}
