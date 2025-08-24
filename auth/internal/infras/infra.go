package infras

import (
	"database/sql"
	"log"

	r "github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/db"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/redis"
)

func Initialize() (*sql.DB, *r.Client) {
	db, err := db.NewConnection()
	if err != nil {
		log.Fatalln("FATAL: database not initialized:", err)
	}
	redis, err := redis.NewConnection()
	if err != nil {
		log.Fatalln("FATAL: redis not initialized:", err)
	}
	return db, redis
}
