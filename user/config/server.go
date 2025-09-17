package config

type serverConfig struct {
	Protocol string
	Host     string
	Port     string
	Timeout  int
}

var serverCfg *serverConfig

func loadServerConfig() {
	serverCfg = &serverConfig{
		Protocol: getEnvWithFallback("SERVER_PROTOCOL", "http"),
		Host:     getEnvWithFallback("SERVER_HOST", "localhost"),
		Port:     getEnvWithFallback("SERVER_PORT", "9001"),
		Timeout:  getNumberEnvWithFallback("SERVER_TIMEOUT", 5),
	}
}

func ServerGetProtocol() (protocol string) {
	return serverCfg.Protocol
}

func ServerGetHost() (host string) {
	return serverCfg.Host
}

func ServerGetPort() (port string) {
	return serverCfg.Port
}

func ServerGetTimeout() (timeout int) {
	return serverCfg.Timeout
}

func ServerGetBaseURL() (url string) {
	return ServerGetHost() + ":" + ServerGetPort()
}
