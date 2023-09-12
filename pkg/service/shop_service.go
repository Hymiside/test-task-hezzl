package service

import (
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
)

type shopService struct {
	repoP *postgres.PostgresRepository
	repoC   *clickhouse.ClickhouseRepository
}

func newShopService(repoP *postgres.PostgresRepository, repoC *clickhouse.ClickhouseRepository) *shopService {
	return &shopService{repoP: repoP, repoC: repoC}
}
