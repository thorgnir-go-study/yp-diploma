package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
)

type Config struct {
	DatabaseDSN string `env:"DATABASE_URI"`
}

func main() {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("configuration failure: failed to parse environment: %w", err)
	}
	flag.StringVar(&cfg.DatabaseDSN, "database_uri", cfg.DatabaseDSN, "Database DSN. Required to be set in CLI or env variable DATABASE_URI")

	flag.Parse()

	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	err = goose.Up(db, "../../migrations")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All migrations done")
}
