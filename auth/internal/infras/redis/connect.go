package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	r "github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/config"
)

func NewConnection() (redis *r.Client, err error) {
	password := config.GetRedisPass()
	if password == "" {
		log.Println("WARNING: connecting to redis with no password")
	}

	client := r.NewClient(&r.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetRedisHost(), config.GetRedisPort()),
		Password: password,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("SUCCESS: connected to redis")

	return client, nil
}
