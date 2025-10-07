package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritchieridanko/apotekly-api/user/config"
	"github.com/ritchieridanko/apotekly-api/user/internal/infrastructure"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/di"
	"github.com/ritchieridanko/apotekly-api/user/internal/server"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatalln("FATAL -> ", err.Error())
	}

	infra, err := infrastructure.Initialize(*cfg)
	if err != nil {
		log.Fatalln("FATAL -> ", err.Error())
	}
	defer infra.Close()

	c := di.NewContainer(cfg, infra)

	s := server.NewHTTPServer(cfg, c.Router().Engine())
	go s.Start()

	// handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Println("HTTP Server -> forced to shutdown:", err.Error())
	}
}
