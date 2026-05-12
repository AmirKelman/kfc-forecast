package config

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Forecast ForecastConfig `mapstructure:"forecast"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ForecastConfig struct {
	GenerationCron string `mapstructure:"generation_cron"`
	HistoryDays    int    `mapstructure:"history_days"`
	DaysAhead      int    `mapstructure:"days_ahead"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
