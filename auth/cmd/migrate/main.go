package main

import (
	"flag"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/db"
)

func main() {
	// Define CLI flags
	up := flag.Bool("up", false, "Run all migrations up")
	down := flag.Int("down", 0, "Rollback N migrations down")
	flag.Parse()

	// Load .env
	config.Initialize()

	// Connect DB
	db, err := db.NewConnection()
	if err != nil {
		log.Fatalln("FATAL: database not initialized:", err)
	}
	defer db.Close()

	if *up {
		if err := database.RunMigrations(db); err != nil {
			log.Fatalln("FATAL: failed to run migrations:", err)
		}
	} else if *down >= 0 {
		if err := database.RollbackMigrations(db, *down); err != nil {
			log.Fatalln("FATAL: failed to rollback migrations:", err)
		}
	} else {
		log.Println("no action specified (use -up, -down, or -down N)")
	}
}
