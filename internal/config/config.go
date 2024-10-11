package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func NewConfig() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &c, nil
}

type Config struct {
	DBUrl  string       `env:"DB_URL"`
	Port   int          `env:"PORT" envDefault:"8080"`
	Local  bool         `env:"LOCAL" envDefault:"false"`
	JWT    JWTConfig    `envPrefix:"JWT_"`
	Admin  AdminConfig  `envPrefix:"ADMIN_"`
	Logger LoggerConfig `envPrefix:"LOG_"`
}

type JWTConfig struct {
	Secret           string `env:"JWT_SECRET"`
	AccessExpiration string `env:"JWT_ACCESS_EXPIRATION"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type AdminConfig struct {
	Username string `env:"ADMIN_USERNAME" validate:"min=3,max=24"`
	Email    string `env:"ADMIN_EMAIL" validate:"email"`
	Password string `env:"ADMIN_PASSWORD" validate:"password"`
}
