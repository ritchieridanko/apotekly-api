package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/ritchieridanko/apotekly-api/user/config"
	"github.com/ritchieridanko/apotekly-api/user/internal/infrastructure/database"
	"github.com/ritchieridanko/apotekly-api/user/internal/infrastructure/storage"
	"github.com/ritchieridanko/apotekly-api/user/internal/infrastructure/tracer"
)

type Infrastructure struct {
	db      *sql.DB
	storage *cloudinary.Cloudinary
	tracer  *tracer.Tracer
}

func Initialize(cfg config.Config) (*Infrastructure, error) {
	db, err := database.NewConnection(
		cfg.Database.DSN,
		cfg.Database.MaxIdleConns,
		cfg.Database.MaxOpenConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		return nil, err
	}

	s, err := storage.NewInstance(cfg.Storage.Bucket, cfg.Storage.APIKey, cfg.Storage.APISecret)
	if err != nil {
		return nil, err
	}

	t, err := tracer.NewProvider(cfg.App.Name, cfg.Tracer.Endpoint)
	if err != nil {
		return nil, err
	}

	return &Infrastructure{db: db, storage: s, tracer: t}, nil
}

func (i *Infrastructure) DB() *sql.DB {
	return i.db
}

func (i *Infrastructure) Storage() *cloudinary.Cloudinary {
	return i.storage
}

func (i *Infrastructure) Tracer() *tracer.Tracer {
	return i.tracer
}

func (i *Infrastructure) Close() error {
	if err := i.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	i.tracer.Cleanup()
	return nil
}
