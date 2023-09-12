package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hymiside/test-task-hezzl/pkg/handler"
	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/clickhouse"
	"github.com/Hymiside/test-task-hezzl/pkg/server"
	"github.com/Hymiside/test-task-hezzl/pkg/service"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Panicf("error .env file not found: %v", err)
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

	dbC, err := clickhouse.NewClickhouseDB(
		ctx, 
		models.ConfigClickhouseRepository{
			Host: os.Getenv("HOST_CLICKHOUSE"),
			Port: os.Getenv("PORT_CLICKHOUSE"),
			Name: os.Getenv("NAME_CLICKHOUSE"),
		})
	if err != nil {
		log.Panicf("error to init clickhouse repository: %v", err)
	}

	repoP := postgres.NewPostgresRepository(dbP)
	repoC := clickhouse.NewClickhouseRepository(dbC)

	services := service.NewService(repoP, repoC)
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