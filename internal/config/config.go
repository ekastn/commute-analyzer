package config

import "github.com/ekastn/commute-analyzer/internal/env"

type Config struct {
	Addr        string
	DatabaseURL string
	ORSAPIKey   string
}

func Load() Config {
	return Config{
		Addr:        env.GetString("SRV_ADDR", ":8080"),
		DatabaseURL: env.GetString("DATABASE_URL", ""),
		ORSAPIKey:   env.GetString("ORS_API_KEY", ""),
	}
}
