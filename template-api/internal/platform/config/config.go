package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds static, deploy-time configuration loaded from the environment.
// This is wiring + secrets: it is read once at startup and is immutable for the
// life of the process. Runtime-mutable settings belong elsewhere (a DB-backed
// store), never here.
type Config struct {
	AppEnv      string   `env:"APP_ENV" envDefault:"development"`
	Port        string   `env:"APP_PORT" envDefault:"8080"`
	DatabaseURL string   `env:"DATABASE_URL,required"`
	CORSOrigins []string `env:"CORS_ORIGINS" envSeparator:","`
}

// Load reads .env (if present) then parses environment variables into Config.
func Load() (*Config, error) {
	// .env is optional; ignore "not found" so prod can rely on real env vars.
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

// IsProduction reports whether the app runs in production mode.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
