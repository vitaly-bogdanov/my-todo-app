package core_http_server

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config — настройки HTTP-сервера из переменных HTTP_*.
// Port берётся из HTTP_PORT; если HTTP_PORT не задан, используется PORT
// (его выставляют PaaS вроде Render).
type Config struct {
	Host               string        `envconfig:"HOST"`
	Port               int           `envconfig:"PORT"`
	Addr               string        `ignored:"true"`
	ShutdownTimeout    time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	CORSAllowedOrigins []string      `envconfig:"CORS_ALLOWED_ORIGINS" required:"true"`
	CORSAllowedMethods []string      `envconfig:"CORS_ALLOWED_METHODS" default:"GET,POST,PATCH,PUT,DELETE,OPTIONS"`
	CORSAllowedHeaders []string      `envconfig:"CORS_ALLOWED_HEADERS" default:"Content-Type,Authorization"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("failed to process http config: %w", err)
	}

	if config.Port == 0 {
		if port := os.Getenv("PORT"); port != "" {
			p, err := strconv.Atoi(port)
			if err != nil {
				return Config{}, fmt.Errorf("invalid PORT: %q: %w", port, err)
			}
			config.Port = p
		}
	}

	if config.Port < 1 || config.Port > 65535 {
		return Config{}, fmt.Errorf("http port must be between 1 and 65535: %d", config.Port)
	}

	if config.Host == "" {
		config.Addr = fmt.Sprintf(":%d", config.Port)
	} else {
		config.Addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("failed to get http config: %w", err)
		panic(err)
	}
	return config
}
