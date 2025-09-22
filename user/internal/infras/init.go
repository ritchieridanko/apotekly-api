package infras

import (
	"database/sql"
	"log"

	c "github.com/cloudinary/cloudinary-go/v2"
	"github.com/ritchieridanko/apotekly-api/user/internal/infras/cloudinary"
	ot "github.com/ritchieridanko/apotekly-api/user/internal/infras/open_telemetry"
	"github.com/ritchieridanko/apotekly-api/user/internal/infras/postgresql"
)

func Initialize() (db *sql.DB, tracer *ot.Tracer, storage *c.Cloudinary) {
	db, err := postgresql.Connect()
	if err != nil {
		log.Fatalln("FATAL -> failed to connect to database:", err)
	}
	tracer, err = ot.Initialize()
	if err != nil {
		log.Fatalln("FATAL -> failed to initialize tracer:", err)
	}
	storage, err = cloudinary.Initialize()
	if err != nil {
		log.Fatalln("FATAL -> failed to initialize cloudinary:", err)
	}
	return db, tracer, storage
}
