package config

type clientConfig struct {
	BaseURL string
}

var clientCfg *clientConfig

func loadClientConfig() {
	clientCfg = &clientConfig{
		BaseURL: getEnv("CLIENT_URL"),
	}
}

func ClientGetBaseURL() (url string) {
	return clientCfg.BaseURL
}
