package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/di"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/validator"
	"github.com/ritchieridanko/apotekly-api/auth/internal/servers"
)

func main() {
	cfg, err := configs.Load("./configs")
	if err != nil {
		log.Fatalln("FATAL -> ", err.Error())
	}

	infra, err := infrastructure.Initialize(cfg)
	if err != nil {
		log.Fatalln("FATAL -> ", err.Error())
	}
	defer infra.Close()

	c := di.NewContainer(cfg, infra)

	if err := validator.RegisterValidators(); err != nil {
		log.Fatalln("FATAL -> ", err.Error())
	}

	s := servers.NewHTTPServer(cfg, c.Router().Engine())
	go func() {
		if err := s.Start(); err != nil {
			log.Fatalln("FATAL -> ", err.Error())
		}
	}()

	// handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Println("FORCED TO SHUTDOWN -> ", err.Error())
	}
}
