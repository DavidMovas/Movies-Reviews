package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &c, nil
}

type Config struct {
	DBUrl      string           `env:"DB_URL"`
	Port       int              `env:"PORT" envDefault:"8000" envRequired:"true"`
	Local      bool             `env:"LOCAL" envDefault:"true"`
	JWT        JWTConfig        `envPrefix:"JWT_"`
	Admin      AdminConfig      `envPrefix:"ADMIN_"`
	Logger     LoggerConfig     `envPrefix:"LOG_"`
	Pagination PaginationConfig `envPrefix:"PAGINATION_"`
}

type JWTConfig struct {
	Secret           string        `env:"SECRET"`
	AccessExpiration time.Duration `env:"ACCESS_EXPIRATION"`
}

type LoggerConfig struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

type AdminConfig struct {
	Username string `env:"USERNAME" validate:"min=3,max=24"`
	Email    string `env:"EMAIL" validate:"email"`
	Password string `env:"PASSWORD" validate:"password"`
}

type PaginationConfig struct {
	DefaultSize int `env:"DEFAULT_SIZE" envDefault:"10"`
	MaxSize     int `env:"MAX_SIZE" envDefault:"20"`
}
