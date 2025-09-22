package config

type storageConfig struct {
	Name      string
	APIKey    string
	APISecret string

	ImageMaxSize int64
}

var storageCfg *storageConfig

func loadStorageConfig() {
	storageCfg = &storageConfig{
		Name:      getEnv("STORAGE_CLOUD_NAME"),
		APIKey:    getEnv("STORAGE_API_KEY"),
		APISecret: getEnv("STORAGE_API_SECRET"),

		ImageMaxSize: int64(getNumberEnvWithFallback("STORAGE_IMAGE_MAX_SIZE", 1)), // fallback: 1 MB
	}
}

func StorageGetName() (name string) {
	return storageCfg.Name
}

func StorageGetAPIKey() (key string) {
	return storageCfg.APIKey
}

func StorageGetAPISecret() (secret string) {
	return storageCfg.APISecret
}

func StorageGetImageMaxSize() (max int64) {
	return storageCfg.ImageMaxSize
}
