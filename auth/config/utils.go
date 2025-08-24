package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv(key string) (value string) {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("FATAL: env %s not found", key)
	}
	return value
}

func GetEnvWithFallback(key, fallback string) (value string) {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("WARNING: env %s not found, using fallback %s", key, fallback)
		return fallback
	}
	return value
}

func GetNumberEnv(key string) (value int) {
	number, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("FATAL: env %s not found", key)
	}
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Fatalf("FATAL: env %s is not a number", key)
	}
	return value
}

func GetNumberEnvWithFallback(key string, fallback int) (value int) {
	number, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("WARNING: env %s not found, using fallback %d", key, fallback)
		return fallback
	}
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Printf("WARNING: env %s is not a number, using fallback %d", key, fallback)
		return fallback
	}
	return value
}
