package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBUrl string    `env:"DB_URL"`
	Port  int       `env:"PORT" envDefault:"8080"`
	JWT   JWTConfig `envPrefix:"JWT_"`
}

type JWTConfig struct {
	Secret           string `env:"JWT_SECRET"`
	AccessExpiration string `env:"JWT_ACCESS_EXPIRATION"`
}

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &c, nil
}
