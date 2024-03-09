package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"os"
)

type Config struct {
	PostgresHost string
	LogLevel     string

	TelegramToken string
}

func LoadConfig(filePath string) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		zap.L().Error("Error loading .env file", zap.Error(err))
		return nil, err
	}

	cfg := Config{} // 👈 new instance of `Config`

	err = env.Parse(&cfg) // 👈 Parse environment variables into `Config`
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}
	cfg.PostgresHost = postgresHost

	// set sensitive data
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return nil, errors.New("TELEGRAM_TOKEN is not set")
	}
	cfg.TelegramToken = telegramToken

	return &cfg, nil
}
