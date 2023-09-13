package service

import (
	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/redis"
)

type shop interface {
	Create(data models.Good) (models.Good, error)
	Update(data models.Good) (models.Good, error)
	Delete(data models.Good) (map[string]string, error)
}

type Service struct {
	Shop shop
}

func NewService(repoP *postgres.PostgresRepository, repoC *clickhouse.ClickhouseRepository, ch *redis.RedisRepository) *Service {
	return &Service{Shop: newShopService(repoP, repoC, ch)}
}