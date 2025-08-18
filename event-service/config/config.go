package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"go-simpler.org/env"
)

type Cfg struct {
	AppPort          string `env:"APP_PORT" default:"8081"`
	RedisURL         string `env:"REDIS_URL"`
	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabasePort     string `env:"DATABASE_PORT"`
	DatabaseUser     string `env:"DATABASE_USER"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabaseName     string `env:"DATABASE_NAME"`
}

func LoadConfig() (*Cfg, error) {

	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := Cfg{}

	if err := env.Load(&cfg, nil); err != nil {
		return nil, err
	}
	return &cfg, nil
}
