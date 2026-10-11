package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

// NewConfig читает Config из переменных окружения.
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("failed to process logger config: %w", err)
	}

	return config, nil
}

// NewConfigMust читает Config из переменных окружения.
// Если прочитать настройки не удалось, NewConfigMust паникует.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("failed to get logger config: %w", err)
		panic(err)
	}
	return config
}
