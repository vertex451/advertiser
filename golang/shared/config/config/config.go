package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
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

func Load() *Config {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %e", err)
	}

	cfg := Config{}
	err = env.Parse(&cfg) // Parse environment variables into `Config`
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	return &cfg
}
