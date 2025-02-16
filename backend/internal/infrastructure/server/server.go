package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/logger"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/router"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	router     *router.Router
	logger     logger.Logger
	cfg        *config.Config
}

func NewServer(
	router *router.Router,
	logger logger.Logger,
	cfg *config.Config,
) *Server {
	return &Server{
		router: router,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	gin.SetMode(s.cfg.Server.Mode)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler:      s.router.Engine(),
		ReadTimeout:  s.cfg.Server.TimeoutRead,
		WriteTimeout: s.cfg.Server.TimeoutWrite,
		IdleTimeout:  s.cfg.Server.TimeoutIdle,
	}

	s.logger.Info("starting server",
		logger.String("address", s.httpServer.Addr),
		logger.String("mode", s.cfg.Server.Mode),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
