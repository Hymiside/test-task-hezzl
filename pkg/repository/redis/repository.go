package redis

import (
	"context"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/go-redis/cache/v9"
)

type redisCh interface {
	GetAll(ctx context.Context) ([]models.Good, error)
	SetItems(ctx context.Context, items []models.Good) error
	DeleteItem(ctx context.Context) error
}

type RedisRepository struct {
	Ch redisCh
}

func NewRedisRepository(ch *cache.Cache) *RedisRepository {
	return &RedisRepository{Ch: newRedisCache(ch)}
}
