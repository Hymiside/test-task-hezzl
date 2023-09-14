package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hymiside/test-task-hezzl/pkg/handler"
	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/natsqueue"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/redis"
	"github.com/Hymiside/test-task-hezzl/pkg/server"
	"github.com/Hymiside/test-task-hezzl/pkg/service"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Panicf("error to load .env file: %v", err)
	}

	dbP, err := postgres.NewPostgresDB(
		ctx,
		models.ConfigPostgresRepository{
			Host:     os.Getenv("HOST_DB"),
			Port:     os.Getenv("PORT_DB"),
			User:     os.Getenv("USER_DB"),
			Password: os.Getenv("PASSWORD_DB"),
			Name:     os.Getenv("NAME_DB"),
		})
	if err != nil {
		log.Panicf("error to init postgres repository: %v", err)
	}

	ch, err := redis.NewRedisDB(ctx, models.ConfigRedis{
		Host: os.Getenv("HOST_REDIS"),
		Port: os.Getenv("PORT_REDIS"),
	})
	if err != nil {
		log.Panicf("error to init redis repository: %v", err)
	}

	nc, err := natsqueue.NewNatsConn(ctx, models.ConfigNats{
		Host: os.Getenv("HOST_NATS"),
		Port: os.Getenv("PORT_NATS"),
	})
	if err != nil {
		log.Panicf("error to init nats queue: %v", err)
	}

	repoR := redis.NewRedisRepository(ch)
	repoP := postgres.NewPostgresRepository(dbP)

	ncQ := natsqueue.NewNatsQueue(nc, repoP)

	services := service.NewService(repoP, repoR, ncQ)
	handlers := handler.NewHandler(services)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-quit:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	srv := server.Server{}
	if err = srv.RunServer(ctx, handlers.InitRoutes(), models.ConfigServer{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT")}); err != nil {
		log.Panicf("failed to run server: %v", err)
	}
}
