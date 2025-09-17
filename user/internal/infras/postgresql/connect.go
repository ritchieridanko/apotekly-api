package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/apotekly-api/user/config"
)

func Connect() (db *sql.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBGetHost(),
		config.DBGetPort(),
		config.DBGetUser(),
		config.DBGetPass(),
		config.DBGetName(),
		config.DBGetSSLMode(),
	)

	db, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.DBGetMaxIdleConns())
	db.SetMaxOpenConns(config.DBGetMaxOpenConns())
	db.SetConnMaxLifetime(time.Duration(config.DBGetConnMaxLifetime()))

	log.Println("SUCCESS -> connected to database")
	return db, nil
}
