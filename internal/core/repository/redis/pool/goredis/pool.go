package core_repository_redis_pool_goredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	core_repository_redis_pool "github.com/vitaly-bogdanov/my-todo-app/internal/core/repository/redis/pool"
)

type Pool struct {
	client *redis.Client
	ttl    time.Duration
	cfg    Config
}

func NewPool(ctx context.Context, cfg Config) (*Pool, error) {
	var options *redis.Options
	if cfg.URL != "" {
		var err error
		options, err = redis.ParseURL(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("parse REDIS_URL: %w", err)
		}
	} else {
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		}
	}

	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Pool{client: client, ttl: cfg.TTL, cfg: cfg}, nil
}

func (p *Pool) Get(ctx context.Context, key string) core_repository_redis_pool.StringCmd {
	return stringCmd{cmd: p.client.Get(ctx, key)}
}

func (p *Pool) Set(ctx context.Context, key string, value any, ttl time.Duration) core_repository_redis_pool.StatusCmd {
	return statusCmd{cmd: p.client.Set(ctx, key, value, ttl)}
}

func (p *Pool) Del(ctx context.Context, keys ...string) core_repository_redis_pool.IntCmd {
	return intCmd{cmd: p.client.Del(ctx, keys...)}
}

func (p *Pool) HGet(ctx context.Context, key string, field string) core_repository_redis_pool.StringCmd {
	return stringCmd{cmd: p.client.HGet(ctx, key, field)}
}

func (p *Pool) HSet(ctx context.Context, key string, values ...any) core_repository_redis_pool.IntCmd {
	return intCmd{cmd: p.client.HSet(ctx, key, values...)}
}

func (p *Pool) Expire(ctx context.Context, key string, ttl time.Duration) core_repository_redis_pool.StatusCmd {
	return statusCmd{cmd: p.client.Expire(ctx, key, ttl)}
}

func (p *Pool) Close() error {
	return p.client.Close()
}

func (p *Pool) TTL() time.Duration {
	return p.ttl
}
