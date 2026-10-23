package core_repository_redis_pool_goredis

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	URL      string        `envconfig:"URL"`
	Host     string        `envconfig:"HOST"`
	Port     int           `envconfig:"PORT" default:"6379"`
	Password string        `envconfig:"PASSWORD"`
	DB       int           `envconfig:"DB" default:"0"`
	TTL      time.Duration `envconfig:"TTL" default:"5m"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("failed to process REDIS environment variables: %w", err)
	}
	if config.URL == "" && config.Host == "" {
		return Config{}, fmt.Errorf("either REDIS_URL or REDIS_HOST is required")
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("failed to create REDIS config: %w", err)
		panic(err)
	}
	return config
}
