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

	loc, err := time.LoadLocation(cfg.Forecast.Timezone)
	if err != nil {
		log.Fatalf("invalid forecast.timezone %q: %v", cfg.Forecast.Timezone, err)
	}

	database, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	forecastSvc := forecast.NewService(database, cfg.Forecast.HistoryDays, cfg.Forecast.DaysAhead, loc)

	sched := scheduler.New(forecastSvc, database, loc)
	if err := sched.Start(cfg.Forecast.GenerationCron); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}

	srv := server.New(database, forecastSvc, cfg.Forecast.AdminToken)
	httpServer := srv.HTTPServer(cfg.Server.Port)

	go func() {
		log.Printf("server listening on :%d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

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
