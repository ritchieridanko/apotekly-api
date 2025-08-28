package infras

import (
	"database/sql"
	"log"

	"github.com/ritchieridanko/apotekly-api/user/internal/infras/db"
)

func Initialize() *sql.DB {
	db, err := db.NewConnection()
	if err != nil {
		log.Fatalln("FATAL: database not initialized:", err)
	}
	return db
}
