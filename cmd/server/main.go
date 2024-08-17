package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dev.theenthusiast.safe-store/internal/api"
	"dev.theenthusiast.safe-store/internal/config"
	"dev.theenthusiast.safe-store/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configurations", "error", err)
	}

	server, err := api.NewServer(cfg, log)
	if err != nil {
		log.Fatal("failed to create server", "error", err)
	}
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
