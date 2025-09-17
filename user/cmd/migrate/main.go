package main

import (
	"flag"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/apotekly-api/user/config"
	"github.com/ritchieridanko/apotekly-api/user/database"
	"github.com/ritchieridanko/apotekly-api/user/internal/infras/postgresql"
)

func main() {
	// define CLI flags
	up := flag.Bool("up", false, "Run all migrations up")
	down := flag.Int("down", 0, "Rollback N migrations down")
	flag.Parse()

	// load .env
	config.Initialize()

	// connect to DB
	db, err := postgresql.Connect()
	if err != nil {
		log.Fatalln("FATAL -> failed to connect to database:", err)
	}
	defer db.Close()

	if *up {
		if err := database.RunMigrations(db); err != nil {
			log.Fatalln("FATAL -> failed to apply migrations:", err)
		}
	} else if *down >= 0 {
		if err := database.RollbackMigrations(db, *down); err != nil {
			log.Fatalln("FATAL -> failed to rollback migrations:", err)
		}
	} else {
		log.Println("WARNING -> no action specified (use -up, -down, or -down N)")
	}
}
