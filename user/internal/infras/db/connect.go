package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/apotekly-api/user/config"
)

func NewConnection() (db *sql.DB, err error) {
	host := config.GetDBHost()
	port := config.GetDBPort()
	user := config.GetDBUser()
	pass := config.GetDBPass()
	name := config.GetDBName()
	mode := config.GetDBSSLMode()
	maxIdleConns := config.GetDBMaxIdleConns()
	maxOpenConns := config.GetDBMaxOpenConns()
	connMaxLifetime := config.GetDBConnMaxLifetime()

	connUrl := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		pass,
		name,
		mode,
	)

	db, err = sql.Open("pgx", connUrl)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime))

	log.Println("SUCCESS: connected to the database")

	return db, nil
}
