package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kfc-forecast/internal/config"
	"kfc-forecast/internal/db"
	"kfc-forecast/internal/forecast"
	"kfc-forecast/internal/scheduler"
	"kfc-forecast/internal/server"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	forecastSvc := forecast.NewService(database, cfg.Forecast.HistoryDays, cfg.Forecast.DaysAhead)

	sched := scheduler.New(forecastSvc, database)
	if err := sched.Start(cfg.Forecast.GenerationCron); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}

	srv := server.New(database, forecastSvc)
	httpServer := srv.HTTPServer(cfg.Server.Port)

	// Start HTTP server in background.
	go func() {
		log.Printf("server listening on :%d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until SIGINT or SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	sched.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("shutdown complete")
}
