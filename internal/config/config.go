package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AllowOrigins         []string `env:"ALLOW_ORIGINS" envSeparator:","`
	LogFormat            string   `env:"LOG_FORMAT"`
	Port                 string   `env:"PORT"`
	JWT                  JWT      `envPrefix:"JWT_"`
	CORS                 CORS     `envPrefix:"CORS_"`
	Database             Database `envPrefix:"DATABASE_"`
	GoogleAppCredentials string   `env:"GOOGLE_APP_CREDENTIALS"`
	UploadSlipBucket     string   `env:"UPLOAD_SLIP_BUCKET"`
}

// @WireSet("Config")
func NewConfig() *Config {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().
			Msg("No .env file found or error loading it. Falling back to system environment variables.")
	}

	config := &Config{}

	// Parse environment variables into the config struct
	if err := env.Parse(config); err != nil {
		log.Panic().
			Err(err).
			Msg("Failed to parse environment variables")
	}

	return config
}
