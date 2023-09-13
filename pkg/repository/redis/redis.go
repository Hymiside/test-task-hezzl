package redis

import (
	"context"
	"fmt"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisDB(ctx context.Context, c models.ConfigRedis) (*cache.Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	ch := cache.New(&cache.Options{
		Redis: rdb,
	})

	go func(ctx context.Context) {
		<-ctx.Done()
		rdb.Close()
	}(ctx)

	return ch, nil
}