package service

import (
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
)

type shopService struct {
	postg *postgres.PostgresRepository
	che   *clickhouse.ClickhouseRepository
}

func newShopService(postg *postgres.PostgresRepository, che *clickhouse.ClickhouseRepository) *shopService {
	return &shopService{postg: postg, che: che}
}
