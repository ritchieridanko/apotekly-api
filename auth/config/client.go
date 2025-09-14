package config

type clientConfig struct {
	baseURL string
}

var clientCfg *clientConfig

func LoadClientConfig() {
	clientCfg = &clientConfig{
		baseURL: GetEnv("CLIENT_URL"),
	}
}

func GetClientBaseURL() (url string) {
	return clientCfg.baseURL
}
