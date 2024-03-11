package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

const (
	EnvIntegration = "integration"
	EnvDev         = "dev"
	EnvProd        = "prod"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`

	PostgresHost string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort string `env:"POSTGRES_PORT" envDefault:"5432"`

	Secrets Secrets
}

type Secrets struct {
	TelegramToken string `env:"TELEGRAM_TOKEN,required"`
}

func LoadConfig(filePath string) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		zap.L().Error("Error loading .env file", zap.Error(err))
		return nil, err
	}

	cfg := Config{}
	err = env.Parse(&cfg) // Parse environment variables into `Config`
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	fmt.Println("### cfg", cfg)

	return &cfg, nil
}
