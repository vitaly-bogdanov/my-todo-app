package core_http_middleware

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type CORSConfig struct {
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" required:"true"`
	AllowedMethods []string `envconfig:"ALLOWED_METHODS" default:"GET,POST,PATCH,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `envconfig:"ALLOWED_HEADERS" default:"Content-Type,Authorization"`
}

func NewCORSConfig() (CORSConfig, error) {
	var config CORSConfig
	if err := envconfig.Process("CORS", &config); err != nil {
		return CORSConfig{}, fmt.Errorf("failed to process cors config: %w", err)
	}

	return config, nil
}

func NewCORSConfigMust() CORSConfig {
	config, err := NewCORSConfig()
	if err != nil {
		panic(fmt.Errorf("failed to get cors config: %w", err))
	}

	return config
}
