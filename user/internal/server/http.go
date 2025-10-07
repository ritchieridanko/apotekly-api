package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ritchieridanko/apotekly-api/user/config"
)

type HTTPServer struct {
	server *http.Server
	cfg    *config.Config
}

func NewHTTPServer(cfg *config.Config, handler http.Handler) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	s := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &HTTPServer{server: s, cfg: cfg}
}

func (s *HTTPServer) Start() error {
	log.Println("HTTP Server -> starting on:", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Println("HTTP Server -> shutting down...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	return nil
}
