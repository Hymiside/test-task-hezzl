package service

import (
	"context"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/natsqueue"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/redis"
)

type service interface {
	Create(ctx context.Context, data models.Good) (models.Good, error)
	Update(ctx context.Context, data models.Good) (models.Good, error)
	Delete(ctx context.Context, data models.Good) (interface{}, error)
	GetAll(ctx context.Context, limit, offset int) (map[string]interface{}, error)
}

type Service struct {
	Shop service
}

func NewService(repoP *postgres.PostgresRepository, ch *redis.RedisRepository, nc *natsqueue.Queue) *Service {
	return &Service{Shop: newShopService(repoP, ch, nc)}
}
