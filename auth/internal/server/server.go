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

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/di"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/mailer"
	"github.com/ritchieridanko/apotekly-api/auth/internal/validators"
)

type App struct {
	router *gin.Engine
	server *http.Server
	db     *sql.DB
	cache  *redis.Client
	mailer mailer.Mailer
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	// initialize configurations
	config.Initialize()

	// initialize infrastructures
	db, cache, mailer, tracer := infras.Initialize()
	a.db = db
	a.cache = cache
	a.mailer = mailer
	defer a.db.Close()
	defer a.cache.Close()
	defer a.mailer.Close()
	defer tracer.Cleanup()

	// initialize dependencies
	router := di.SetupDependencies(a.db, a.cache, a.mailer)
	a.router = router

	// initialize validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validators.Initialize(v); err != nil {
			log.Fatalln("FATAL -> failed to initialize validators:", err)
		}
	}

	// create HTTP server
	a.server = &http.Server{
		Addr:    config.ServerGetBaseURL(),
		Handler: a.router,
	}

	// start server
	go func() {
		log.Println("SUCCESS -> running server on:", config.ServerGetBaseURL())
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("FATAL -> failed to start server:", err)
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
		log.Println("STOPPED -> server forced to shutdown:", err)
	}
}
