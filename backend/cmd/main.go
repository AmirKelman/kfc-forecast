package main

import (
	"log"
	"os"

	"kfc-forecast/internal/config"
	"kfc-forecast/internal/db"
	"kfc-forecast/internal/forecast"
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

	srv := server.New(database, forecastSvc)
	log.Printf("server starting on port %d", cfg.Server.Port)
	if err := srv.Run(cfg.Server.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
