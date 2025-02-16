package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/container"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create dependency container
	cont := container.NewContainer(cfg)
	defer cont.Cleanup()

	// Get logger
	log := cont.GetLogger()

	// Set up and start server
	server := cont.GetHTTPServer()

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Error("server error",
				logger.Error(err), // Using the proper logger.Error field constructor
			)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown",
			logger.Error(err), // Using the proper logger.Error field constructor
		)
		os.Exit(1)
	}

	log.Info("server stopped")
}
