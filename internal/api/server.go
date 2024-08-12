package api

import (
	"context"
	"fmt"
	"net/http"

	"dev.theenthusiast.safe-store/internal/config"
	"dev.theenthusiast.safe-store/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router *httprouter.Router
	logger logger.Logger
	config *config.Config
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	return &Server{
		router: httprouter.New(),
		logger: log,
		config: cfg,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", "port", s.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// In a real-world scenario, you might want to implement proper shutdown logic here
	return nil
}
