package config

type cacheConfig struct {
	host       string
	port       string
	pass       string
	maxRetries int
	baseDelay  int
}

var cacheCfg *cacheConfig

func LoadCacheConfig() {
	cacheCfg = &cacheConfig{
		host:       GetEnv("REDIS_HOST"),
		port:       GetEnv("REDIS_PORT"),
		pass:       GetEnv("REDIS_PASS"),
		maxRetries: GetNumberEnvWithFallback("REDIS_MAX_RETRIES", 3),
		baseDelay:  GetNumberEnvWithFallback("REDIS_BASE_DELAY", 100),
	}
}

func GetCacheHost() (host string) {
	return cacheCfg.host
}

func GetCachePort() (port string) {
	return cacheCfg.port
}

func GetCachePass() (pass string) {
	return cacheCfg.pass
}

func GetCacheMaxRetries() (max int) {
	return cacheCfg.maxRetries
}

func GetCacheBaseDelay() (delay int) {
	return cacheCfg.baseDelay
}
