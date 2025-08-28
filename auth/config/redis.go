package config

type redisConfig struct {
	host       string
	port       string
	pass       string
	maxRetries int
	baseDelay  int
}

var redisCfg *redisConfig

func LoadRedisConfig() {
	redisCfg = &redisConfig{
		host:       GetEnv("REDIS_HOST"),
		port:       GetEnv("REDIS_PORT"),
		pass:       GetEnv("REDIS_PASS"),
		maxRetries: GetNumberEnvWithFallback("REDIS_MAX_RETRIES", 3),
		baseDelay:  GetNumberEnvWithFallback("REDIS_BASE_DELAY", 100),
	}
}

func GetRedisHost() (host string) {
	return redisCfg.host
}

func GetRedisPort() (port string) {
	return redisCfg.port
}

func GetRedisPass() (pass string) {
	return redisCfg.pass
}

func GetRedisMaxRetries() (max int) {
	return redisCfg.maxRetries
}

func GetRedisBaseDelay() (delay int) {
	return redisCfg.baseDelay
}
