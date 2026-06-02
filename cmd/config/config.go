package config

import "github.com/erickgreco/indoorgrid-system/pkg/env"

type Config struct {
	Port    string
	ApiURL  string
	DBURL   string
	GinMode string
}

func Load() Config {
	return Config{
		Port:    env.GetString("PORT", ":8080"),
		DBURL:   env.GetString("DB_URL", "postgres://admin:adminpassword@localhost:5433/indoorgrid-system?sslmode=disable"),
		GinMode: env.GetString("GIN_MODE", "debug"),
	}
}
