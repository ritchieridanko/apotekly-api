package config

import (
	"log"
	"os"
	"strconv"
)

func getEnv(key string) (value string) {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("FATAL -> env %s not found\n", key)
	}
	return value
}

func getEnvWithFallback(key, fallback string) (value string) {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("WARNING -> env %s not found, using fallback %s\n", key, fallback)
		return fallback
	}
	return value
}

func getNumberEnv(key string) (value int) {
	number, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("FATAL -> env %s not found\n", key)
	}
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Fatalf("FATAL -> env %s not a number\n", key)
	}
	return value
}

func getNumberEnvWithFallback(key string, fallback int) (value int) {
	number, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("WARNING -> env %s not found, using fallback %d\n", key, fallback)
		return fallback
	}
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Printf("WARNING -> env %s not a number, using fallback %d\n", key, fallback)
		return fallback
	}
	return value
}
