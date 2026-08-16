package config

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Rakaa503/AviGo/internal/env"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DatabaseURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING: .env file not found")
	}

	return &Config{
		AppName: env.Get(
			"APP_NAME",
			"AVIGO",
		),

		AppEnv: env.Get(
			"APP_ENV",
			"development",
		),

		AppPort: env.Get(
			"APP_PORT",
			"8080",
		),

		DatabaseURL: env.Required(
			"DATABASE_URL",
		),
	}
}
