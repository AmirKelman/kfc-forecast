package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kfc-forecast/internal/config"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	var (
		database *gorm.DB
		err      error
	)

	for attempt := 1; attempt <= 10; attempt++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err == nil {
			sqlDB, pingErr := database.DB()
			if pingErr == nil {
				if pingErr = sqlDB.Ping(); pingErr == nil {
					log.Println("database connected")
					return database, nil
				}
			}
			err = pingErr
		}
		log.Printf("db connect attempt %d/10 failed: %v — retrying in 2s", attempt, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after 10 attempts: %w", err)
}
