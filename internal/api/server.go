package api

import (
	"context"
	"fmt"
	"net/http"

	"dev.theenthusiast.safe-store/internal/config"
	"dev.theenthusiast.safe-store/internal/middleware"
	"dev.theenthusiast.safe-store/internal/storage"
	"dev.theenthusiast.safe-store/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router   *httprouter.Router
	logger   logger.Logger
	config   *config.Config
	r2Client *storage.R2Client
}

func NewServer(cfg *config.Config, log logger.Logger) (*Server, error) {
	r2Client, err := storage.NewR2Client(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create R2 client: %w", err)
	}

	s := &Server{
		router:   httprouter.New(),
		logger:   log,
		config:   cfg,
		r2Client: r2Client,
	}

	// Apply CORS middleware to all routes
	s.router.GlobalOPTIONS = middleware.HandleCORS()

	return s, nil
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", "port", s.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// In a real-world scenario, you might want to implement proper shutdown logic here
	return nil
}
