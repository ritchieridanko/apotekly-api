package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ritchieridanko/apotekly-api/auth/config"
)

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		config.GetDBName(),
		driver,
	)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("migrations done successfully")
	return nil
}

func RollbackMigrations(db *sql.DB, steps int) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		config.GetDBName(),
		driver,
	)
	if err != nil {
		return err
	}

	if steps == 0 {
		// Roll back everything
		err = migration.Down()
	} else {
		err = migration.Steps(-steps)
	}

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("migrations rolled back successfully")
	return nil
}
