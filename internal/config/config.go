package config

import (
	"flag"
)

type Config struct {
	StorageType  string
	PostgresConn string
	HTTPAddr     string
	BaseURL      string
}

func New() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.StorageType, "storage", "memory", "storage type: memory or postgres")
	flag.StringVar(&cfg.PostgresConn, "pg-conn", "postgres://postgres:password@localhost:5432/shortener?sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&cfg.HTTPAddr, "addr", ":8080", "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "base-url", "http://localhost:8080", "base URL for short links")

	flag.Parse()

	return cfg
}