package core_repository_postgres_pool_pgx

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     int           `envconfig:"PORT" default:"5432"`
	User     string        `envconfig:"USER" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DB" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
	// SSLMode — необязательный sslmode для подключения (disable, require, verify-full ...).
	// По умолчанию disable: так работает локальный docker-compose и внутренняя сеть Render.
	SSLMode string `envconfig:"SSLMODE" default:"disable"`
}

// URL возвращает строку подключения вида scheme://user:password@host:port/db?sslmode=...
// Логин и пароль экранируются, поэтому спецсимволы в них не ломают строку.
func (c Config) URL(scheme string) string {
	u := url.URL{
		Scheme: scheme,
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   "/" + c.Database,
	}
	q := url.Values{}
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return Config{}, fmt.Errorf("failed to process POSTGRES environment variables: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("failed to create POSTGRES config: %w", err)
		panic(err)
	}
	return config
}
