package service

import (
	"context"
	"fmt"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/redis"
)

type shopService struct {
	repoP   *postgres.PostgresRepository
	repoC   *clickhouse.ClickhouseRepository
	ch 		*redis.RedisRepository
}

func newShopService(repoP *postgres.PostgresRepository, repoC *clickhouse.ClickhouseRepository, ch *redis.RedisRepository) *shopService {
	return &shopService{repoP: repoP, repoC: repoC, ch: ch}
}

func (s *shopService) Create(data models.Good) (models.Good, error) {
	res, err := s.repoP.Shop.Create(data)
	if err != nil {
		return models.Good{}, err
	}
	err = s.updateRedis()
	if err != nil {
		fmt.Printf("error to update redis: %v", err)
	}
	return res, nil
}

func (s *shopService) Update(data models.Good) (models.Good, error) {
	res, err := s.repoP.Shop.Update(data)
	if err != nil {
		return models.Good{}, err
	}
	err = s.updateRedis()
	if err != nil {
		fmt.Printf("error to update redis: %v", err)
	}
	return res, nil
}

func (s *shopService) Delete(data models.Good) (map[string]string, error) {
	res, err := s.repoP.Shop.Delete(data)
	if err != nil {
		return nil, err
	}
	err = s.updateRedis()
	if err != nil {
		fmt.Printf("error to update redis: %v", err)
	}
	return res, nil
}

func (s *shopService) updateRedis() error {
	ctx := context.Background()

	res, err := s.repoP.Shop.GetAll()
	if err != nil {
		return err
	}

	err = s.ch.Ch.DeleteItem(ctx)
	if err != nil {
		return err
	}

	err = s.ch.Ch.SetItems(ctx, res)
	if err != nil {
		return err
	}
	return nil
}