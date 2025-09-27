package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/ritchieridanko/apotekly-api/pharmacy/config"
)

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		config.DBGetName(),
		driver,
	)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("SUCCESS -> applied database migrations")
	return nil
}

func RollbackMigrations(db *sql.DB, steps int) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		config.DBGetName(),
		driver,
	)
	if err != nil {
		return err
	}

	if steps == 0 {
		// rollback all migrations
		err = migration.Down()
	} else {
		err = migration.Steps(-steps)
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("SUCCESS -> rolled back database migrations")
	return nil
}
