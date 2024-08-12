package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"dev.theenthusiast.safe-store/internal/api"
	"dev.theenthusiast.safe-store/internal/config"
	"dev.theenthusiast.safe-store/pkg/logger"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configurations", "error", err)
	}

	server := api.NewServer(cfg, log)
	server.SetupRoutes()

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server is shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")

}
