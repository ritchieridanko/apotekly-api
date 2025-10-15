package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
)

func NewConnection(cfg *configs.Cache) (cache *redis.Client, err error) {
	if cfg.Pass == "" {
		log.Println("WARNING -> connecting to cache without password")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	cache = redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: cfg.Pass,
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cache.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping cache: %w", err)
	}

	log.Println("✅ connected to cache")
	return cache, nil
}
