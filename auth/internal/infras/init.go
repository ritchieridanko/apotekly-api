package infras

import (
	"database/sql"
	"log"

	r "github.com/redis/go-redis/v9"
	m "github.com/ritchieridanko/apotekly-api/auth/internal/infras/mailer"
	ot "github.com/ritchieridanko/apotekly-api/auth/internal/infras/open_telemetry"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/postgresql"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/redis"
)

func Initialize() (db *sql.DB, cache *r.Client, mailer m.Mailer, tracer *ot.Tracer) {
	db, err := postgresql.Connect()
	if err != nil {
		log.Fatalln("FATAL -> failed to connect to database:", err.Error())
	}
	cache, err = redis.Connect()
	if err != nil {
		log.Fatalln("FATAL -> failed to connect to redis:", err.Error())
	}
	mailer = m.NewMailer()
	tracer, err = ot.Initialize()
	if err != nil {
		log.Fatalln("FATAL -> failed to initialize tracer:", err.Error())
	}
	return db, cache, mailer, tracer
}
