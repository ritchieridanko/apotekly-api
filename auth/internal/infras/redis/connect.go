package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/config"
)

func Connect() (cache *redis.Client, err error) {
	password := config.CacheGetPass()
	if password == "" {
		log.Println("WARNING -> connecting to redis without password")
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     config.CacheGetHost() + ":" + config.CacheGetPort(),
			Password: password,
		},
	)

	// test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("SUCCESS -> connected to redis")
	return client, nil
}
