package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Initialize() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("FATAL: unable to load .env file:", err)
	}

	LoadAppConfig()
	LoadAuthConfig()
	LoadDBConfig()
	LoadServerConfig()
}
