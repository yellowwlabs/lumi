package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"lumi.yellowlabs.space/internal/config"
	"lumi.yellowlabs.space/internal/server"
	"lumi.yellowlabs.space/internal/services"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := services.InitPostgres(); err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer services.ClosePostgres()

	if err := services.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := services.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer services.CloseRedis()

	srv := server.NewServer()
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
