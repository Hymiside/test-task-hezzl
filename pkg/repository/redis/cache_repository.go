package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/go-redis/cache/v9"
)

type redisCache struct {
	ch *cache.Cache
}

var TTL = time.Minute

func newRedisCache(ch *cache.Cache) *redisCache {
	return &redisCache{ch: ch}
}

func (r *redisCache) GetAll(ctx context.Context) ([]models.Good, error) {
	var items []models.Good
	if err := r.ch.Get(ctx, "DataItems", &items); err != nil {
		return nil, fmt.Errorf("error to get items from redis: %w", err)
	}
	if items == nil {
		return nil, errors.New("not found items in redis")
	}
	return items, nil
}

func (r *redisCache) SetItems(ctx context.Context, items []models.Good) error {
	if err := r.ch.Set(&cache.Item{
		Ctx:   ctx,
		Key:   "DataItems",
		Value: items,
		TTL:   TTL,
	}); err != nil {
		return fmt.Errorf("error to set items in redis: %w", err)
	}
	return nil
}

func (r *redisCache) DeleteItem(ctx context.Context) error {
	if err := r.ch.Delete(ctx, "DataItems"); err != nil {
		return fmt.Errorf("error to delete items in redis: %w", err)
	}
	return nil
}
