package infras

import (
	"database/sql"
	"log"

	c "github.com/cloudinary/cloudinary-go/v2"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/infras/cloudinary"
	ot "github.com/ritchieridanko/apotekly-api/pharmacy/internal/infras/open_telemetry"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/infras/postgresql"
)

func Initialize() (db *sql.DB, tracer *ot.Tracer, storage *c.Cloudinary) {
	db, err := postgresql.Connect()
	if err != nil {
		log.Fatalln("FATAL -> failed to connect to database:", err.Error())
	}
	tracer, err = ot.Initialize()
	if err != nil {
		log.Fatalln("FATAL -> failed to initialize tracer:", err.Error())
	}
	storage, err = cloudinary.Initialize()
	if err != nil {
		log.Fatalln("FATAL -> failed to initialize cloudinary:", err.Error())
	}
	return db, tracer, storage
}
