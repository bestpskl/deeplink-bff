package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type appConfig struct {
	Host string `envconfig:"APP_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"APP_PORT" default:"4000"`
}

type config struct {
	Environment string `envconfig:"ENV" default:"dev"`
	App         appConfig
}

var (
	cfg  config
	once sync.Once
)

func Load() {
	once.Do(func() {
		_ = godotenv.Load()
		envconfig.MustProcess("", &cfg)
	})
}

func Get() config {
	return cfg
}

func (c config) IsDevelop() bool {
	return c.Environment == "dev"
}
