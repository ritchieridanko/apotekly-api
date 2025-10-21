package main

import (
	"flag"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure/database"
)

func main() {
	up := flag.Bool("up", false, "Run all migrations up")
	down := flag.Int("down", 0, "Rollback N migrations down")
	flag.Parse()

	cfg, err := configs.Load("./configs")
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}

	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
	defer db.Close()

	migrator, err := database.NewMigrator(db, "migrations", cfg.Database.Name)
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
	defer migrator.Close()

	if *up {
		if err := migrator.Up(); err != nil {
			log.Fatalln("FATAL ->", err.Error())
		}
	} else if *down >= 0 {
		if err := migrator.Down(*down); err != nil {
			log.Fatalln("FATAL ->", err.Error())
		}
	} else {
		log.Fatalln("FATAL -> failed to run migrations: no action specified (use -up, -down, or -down N)")
	}
}
