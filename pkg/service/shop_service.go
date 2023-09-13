package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/natsqueue"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/redis"
)

type shopService struct {
	repoP *postgres.PostgresRepository
	ch    *redis.RedisRepository
	nc 	  *natsqueue.Queue
}

func newShopService(repoP *postgres.PostgresRepository, ch *redis.RedisRepository, nc *natsqueue.Queue) *shopService {
	go nc.Q.Sub()
	return &shopService{repoP: repoP, ch: ch, nc: nc}
}

func (s *shopService) Create(ctx context.Context, data models.Good) (models.Good, error) {
	res, err := s.repoP.Shop.Create(data)
	if err != nil {
		return models.Good{}, err
	}
	if err = s.updateRedis(ctx); err != nil {
		fmt.Printf("error to update redis: %v", err)
	}

	if err = s.nc.Q.Pub(res); err != nil {
		fmt.Printf("error to pub log (create): %v", err)
	}

	return res, nil
}

func (s *shopService) Update(ctx context.Context, data models.Good) (models.Good, error) {
	res, err := s.repoP.Shop.Update(data)
	if err != nil {
		return models.Good{}, err
	}
	err = s.updateRedis(ctx)
	if err != nil {
		fmt.Printf("error to update redis: %v", err)
	}

	if err = s.nc.Q.Pub(res); err != nil {
		fmt.Printf("error to pub log (update): %v", err)
	}

	return res, nil
}

func (s *shopService) Delete(ctx context.Context, data models.Good) (interface{}, error) {
	res, err := s.repoP.Shop.Delete(data)
	if err != nil {
		return nil, err
	}
	err = s.updateRedis(ctx)
	if err != nil {
		fmt.Printf("error to update redis: %v", err)
	}

	if err = s.nc.Q.Pub(res); err != nil {
		fmt.Printf("error to pub log (update): %v", err)
	}

	response := map[string]string {
		"id": strconv.Itoa(res.Id),
		"project_id": strconv.Itoa(res.ProjectId),
		"removed": strconv.FormatBool(true),
	}

	return response, nil
}

func (s *shopService) GetAll(ctx context.Context, limit, offset int) (map[string]interface{}, error) {
	goods, err := s.ch.Ch.GetAll(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	if goods == nil {
		goods, err = s.repoP.Shop.GetAll(0, 0)
		if err != nil {
			return nil, err
		}
		err = s.ch.Ch.SetItems(ctx, goods)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		if len(goods) < offset {
			goods = []models.Good{}
		} else if len(goods) < offset+limit {
			goods = goods[offset:]
		} else {
			goods = goods[offset : offset+limit]
		}
	}

	inc := 0
	for i := range goods {
		if goods[i].Removed {
			inc++
		}
	}

	meta_data := map[string]int{
		"total":   len(goods),
		"removed": inc,
		"limit":   limit,
		"offset":  offset,
	}

	res_goods := []models.Good{}
	for i := range goods {
		if !goods[i].Removed {
			res_goods = append(res_goods, goods[i])
		}
	}

	res := map[string]interface{}{
		"meta":  meta_data,
		"goods": res_goods,
	}
	return res, nil

}

func (s *shopService) updateRedis(ctx context.Context) error {
	res, err := s.repoP.Shop.GetAll(0, 0)
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
