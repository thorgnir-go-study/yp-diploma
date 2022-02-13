package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS" envDefault:":8081"`
	DatabaseDSN   string `env:"DATABASE_URI"`
	AccrualURL    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret     string `env:"JWT_SECRET" envDefault:"very very secret key"`
}

func GetConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("configuration failure: failed to parse environment: %w", err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address. If not set in CLI or env variable RUN_ADDRESS defaults to ':8081'")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "Database DSN. Required to be set in CLI or env variable DATABASE_URI")
	flag.StringVar(&cfg.AccrualURL, "r", cfg.AccrualURL, "Accrual system address. Required to be set in CLI or env variable ACCRUAL_SYSTEM_ADDRESS")
	flag.StringVar(&cfg.JWTSecret, "jwt_secret", cfg.JWTSecret, "JWT Secret. can be set via CLI or env variable JWT_SECRET")

	flag.Parse()

	if cfg.DatabaseDSN == "" {
		return nil, errors.New("configuration failure: database dsn not set")
	}

	if cfg.AccrualURL == "" {
		return nil, errors.New("configuration failure: accrual system address not set")
	}

	return cfg, nil
}
