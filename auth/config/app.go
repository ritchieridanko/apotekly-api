package config

type appConfig struct {
	Name        string
	Version     string
	Description string
}

var appCfg *appConfig

func loadAppConfig() {
	appCfg = &appConfig{
		Name:        getEnv("APP_NAME"),
		Version:     getEnv("APP_VERSION"),
		Description: getEnv("APP_DESCRIPTION"),
	}
}

func AppGetName() (name string) {
	return appCfg.Name
}

func AppGetVersion() (version string) {
	return appCfg.Version
}

func AppGetDescription() (description string) {
	return appCfg.Description
}
