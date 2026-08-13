package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"gorm.io/gorm"

	"dailyq-api/internal/config"
)

// Server owns the HTTP listener and its lifecycle.
type Server struct {
	cfg    *config.Config
	http   *http.Server
	logger *slog.Logger
}

// New builds a Server with routes wired to db.
func New(cfg *config.Config, db *gorm.DB, log *slog.Logger) *Server {
	handler := NewRouter(cfg, db, log)

	addr := net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port))

	return &Server{
		cfg:    cfg,
		logger: log,
		http: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
	}
}

// Run starts the server and blocks until SIGINT/SIGTERM, then shuts down
// gracefully.
func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("http server listening", "addr", s.http.Addr, "mode", s.cfg.Server.Mode)
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	s.logger.Info("server stopped")
	return nil
}
