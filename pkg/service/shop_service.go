package service

import (
	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
)

type shopService struct {
	repoP *postgres.PostgresRepository
	repoC   *clickhouse.ClickhouseRepository
}

func newShopService(repoP *postgres.PostgresRepository, repoC *clickhouse.ClickhouseRepository) *shopService {
	return &shopService{repoP: repoP, repoC: repoC}
}

func (s *shopService) Create(data models.Good) (models.Good, error) {
	res, err := s.repoP.Shop.Create(data)
	if err != nil {
		return models.Good{}, err
	}
	return res, nil
}

func (s *shopService) Update(data models.Good) (models.Good, error) {
	res, err := s.repoP.Shop.Update(data)
	if err != nil {
		return models.Good{}, err
	}
	return res, nil
}