package service


import (
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
)

type shop interface {}

type Service struct {
	sh shop
}

func NewService(repoP *postgres.PostgresRepository, repoC *clickhouse.ClickhouseRepository) *Service {
	return &Service{sh: newShopService(repoP, repoC)}
}