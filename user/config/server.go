package config

type serverConfig struct {
	protocol string
	host     string
	port     string
	timeout  int
}

var serverCfg *serverConfig

func LoadServerConfig() {
	serverCfg = &serverConfig{
		protocol: GetEnvWithFallback("SERVER_PROTOCOL", "http"),
		host:     GetEnvWithFallback("SERVER_HOST", "localhost"),
		port:     GetEnvWithFallback("SERVER_PORT", "9001"),
		timeout:  GetNumberEnvWithFallback("GRACEFUL_TIMEOUT", 5),
	}
}

func GetServerProtocol() (protocol string) {
	return serverCfg.protocol
}

func GetServerHost() (host string) {
	return serverCfg.host
}

func GetServerPort() (port string) {
	return serverCfg.port
}

func GetServerTimeout() (timeout int) {
	return serverCfg.timeout
}

func GetServerBaseURL() (baseURL string) {
	return serverCfg.host + ":" + serverCfg.port
}
