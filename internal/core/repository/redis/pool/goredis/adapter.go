package core_repository_redis_pool_goredis

import (
	"errors"

	"github.com/redis/go-redis/v9"
	core_repository_redis_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool"
)

func mapError(err error) error {
	if errors.Is(err, redis.Nil) {
		return core_repository_redis_pool.NotFound
	}
	return err
}

type stringCmd struct {
	cmd *redis.StringCmd
}

func (c stringCmd) Result() (string, error) {
	value, err := c.cmd.Result()
	if err != nil {
		return "", mapError(err)
	}
	return value, nil
}

type statusCmd struct {
	cmd interface{ Err() error }
}

func (c statusCmd) Err() error {
	return mapError(c.cmd.Err())
}

type intCmd struct {
	cmd *redis.IntCmd
}

func (c intCmd) Result() (int64, error) {
	value, err := c.cmd.Result()
	if err != nil {
		return 0, mapError(err)
	}
	return value, nil
}
