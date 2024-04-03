package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	EnvIntegration = "integration"
	EnvDev         = "dev"
	EnvStage       = "stage"
	EnvProd        = "prod"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"dev"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"DEBUG"`

	DataDir string `env:"DATA_DIR" envDefault:"/Users/vertex451/tg-bot-data"`

	DbHost string `env:"DB_HOST" envDefault:"localhost"`
	DbPort string `env:"DB_PORT" envDefault:"5432"`
	DbUser string `env:"DB_USER" envDefault:"postgres"`
	DbName string `env:"DB_NAME" envDefault:"postgres"`

	Secrets Secrets
}

type Secrets struct {
	AgencyTgToken string `env:"AGENCY_TG_TOKEN,required"`
	OwnerTgToken  string `env:"OWNER_TG_TOKEN,required"`

	DbPassword string `env:"DB_PASSWORD,required"`

	//AwsAccessKey string `env:"AWS_ACCESS_KEY_ID,required"`
	//AwsSecret    string `env:"AWS_SECRET_ACCESS_KEY,required"`
}

func Load() *Config {
	// load .env file into env variables
	err := godotenv.Load()
	if err != nil {
		zap.L().Warn("no .env file found, using default values")
	}

	// parse environment variables into struct
	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		zap.L().Warn("failed to parse environment variables", zap.Error(err))
	}

	return &cfg
}
