package config

type appConfig struct {
	appName        string
	appVersion     string
	appDescription string
}

var appCfg *appConfig

func LoadAppConfig() {
	appCfg = &appConfig{
		appName:        GetEnv("APP_NAME"),
		appVersion:     GetEnv("APP_VERSION"),
		appDescription: GetEnv("APP_DESCRIPTION"),
	}
}

func GetAppName() (name string) {
	return appCfg.appName
}

func GetAppVersion() (version string) {
	return appCfg.appVersion
}

func GetAppDescription() (description string) {
	return appCfg.appDescription
}
