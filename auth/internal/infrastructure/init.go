package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/tracer"
)

type Infrastructure struct {
	cache  *redis.Client
	db     *sql.DB
	tracer *tracer.Tracer
}

func Initialize(cfg *configs.Config) (*Infrastructure, error) {
	c, err := cache.NewConnection(&cfg.Cache)
	if err != nil {
		return nil, err
	}

	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		return nil, err
	}

	t, err := tracer.NewProvider(cfg)
	if err != nil {
		return nil, err
	}

	return &Infrastructure{cache: c, db: db, tracer: t}, nil
}

func (i *Infrastructure) Cache() *redis.Client {
	return i.cache
}

func (i *Infrastructure) DB() *sql.DB {
	return i.db
}

func (i *Infrastructure) Tracer() *tracer.Tracer {
	return i.tracer
}

func (i *Infrastructure) Close() error {
	if err := i.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache connection: %w", err)
	}
	if err := i.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	i.tracer.Cleanup()
	return nil
}
