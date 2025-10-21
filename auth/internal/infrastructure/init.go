package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/broker"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/logger"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/tracer"
	"go.uber.org/zap"
)

type Infrastructure struct {
	cache  *redis.Client
	db     *sql.DB
	tracer *tracer.Tracer
	logger *zap.Logger
	broker *broker.Broker
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

	l := logger.NewProvider(&cfg.App)
	b := broker.NewClient(&cfg.Broker)

	return &Infrastructure{cache: c, db: db, tracer: t, logger: l, broker: b}, nil
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

func (i *Infrastructure) Logger() *zap.Logger {
	return i.logger
}

func (i *Infrastructure) Broker() *broker.Broker {
	return i.broker
}

func (i *Infrastructure) Close() error {
	if err := i.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache connection: %w", err)
	}
	if err := i.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to flush buffered log entries: %w", err)
	}
	if err := i.broker.Close(); err != nil {
		return fmt.Errorf("failed to close broker: %w", err)
	}

	i.tracer.Cleanup()
	return nil
}
