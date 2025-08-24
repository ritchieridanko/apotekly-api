package cmd

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
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/di"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras"
	"github.com/ritchieridanko/apotekly-api/auth/internal/validators"
)

type App struct {
	router *gin.Engine
	server *http.Server
	db     *sql.DB
}

func (a *App) Run() {
	// Initialize .env configurations
	config.Initialize()

	// Initialize database
	db := infras.Initialize()
	a.db = db
	defer a.db.Close()

	// Initialize dependency injections
	router := di.SetupDependencies(a.db)
	a.router = router

	// Initialize validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validators.Initialize(v); err != nil {
			log.Fatalln("FATAL: failed to initialize validators:", err)
		}
	}

	// Create HTTP server
	a.server = &http.Server{
		Addr:    config.GetServerBaseURL(),
		Handler: a.router,
	}

	// Start server
	go func() {
		log.Println("Starting server on", config.GetServerBaseURL())
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("FATAL: Failed to start server:", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GetServerTimeout())*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}
}
