package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Initialize() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("FATAL -> failed to load .env file:", err.Error())
	}

	loadAppConfig()
	loadAuthConfig()
	loadCacheConfig()
	loadClientConfig()
	loadDBConfig()
	loadMailerConfig()
	loadOAuthConfig()
	loadServerConfig()
	loadTracerConfig()
}
