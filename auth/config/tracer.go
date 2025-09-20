package config

type tracerConfig struct {
	Endpoint string
}

var tracerCfg *tracerConfig

func loadTracerConfig() {
	tracerCfg = &tracerConfig{
		Endpoint: getEnvWithFallback("TRACER_ENDPOINT", "localhost:4318"),
	}
}

func TracerGetEndpoint() (endpoint string) {
	return tracerCfg.Endpoint
}
