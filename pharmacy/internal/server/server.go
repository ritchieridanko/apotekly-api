package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/config"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/di"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/infras"
)

type App struct {
	router  *gin.Engine
	server  *http.Server
	db      *sql.DB
	storage *cloudinary.Cloudinary
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	// initialize configurations
	config.Initialize()

	// initialize infrastructures
	db, tracer, storage := infras.Initialize()
	a.db = db
	a.storage = storage
	defer a.db.Close()
	defer tracer.Cleanup()

	// initialize dependencies
	router := di.SetupDependencies(a.db, a.storage)
	a.router = router

	// create HTTP server
	a.server = &http.Server{
		Addr:    config.ServerGetBaseURL(),
		Handler: a.router,
	}

	// start server
	go func() {
		log.Println("SUCCESS -> running server on:", config.ServerGetBaseURL())
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("FATAL -> failed to start server:", err.Error())
		}
	}()

	// handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("SHUTTING DOWN SERVER...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ServerGetTimeout())*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Println("STOPPED -> server forced to shutdown:", err.Error())
	}
}
