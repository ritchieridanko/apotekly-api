package infras

import (
	"database/sql"
	"log"

	ot "github.com/ritchieridanko/apotekly-api/user/internal/infras/open_telemetry"
	"github.com/ritchieridanko/apotekly-api/user/internal/infras/postgresql"
)

func Initialize() (db *sql.DB, tracer *ot.Tracer) {
	db, err := postgresql.Connect()
	if err != nil {
		log.Fatalln("FATAL -> failed to connect to database:", err)
	}
	tracer, err = ot.Initialize()
	if err != nil {
		log.Fatalln("FATAL -> failed to initialize tracer:", err)
	}
	return db, tracer
}
