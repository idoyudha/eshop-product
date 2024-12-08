package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/config"
	v1 "github.com/idoyudha/eshop-product/internal/controller/http/v1"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/internal/usecase/repo"
	"github.com/idoyudha/eshop-product/pkg/dynamodb"
	"github.com/idoyudha/eshop-product/pkg/httpserver"
	"github.com/idoyudha/eshop-product/pkg/logger"
	"github.com/idoyudha/eshop-product/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	dynamoDB, err := dynamodb.NewDynamoDB(&cfg.DynamoDB)
	if err != nil {
		l.Fatal("app - Run - dynamodb.NewDynamoDB: ", err)
	}

	redisClient := redis.NewRedis(cfg.Redis)

	productUsecase := usecase.New(
		repo.NewProductRepo(dynamoDB),
		repo.NewCategoryDynamoRepo(dynamoDB),
		repo.NewCategoryRedisRepo(redisClient),
	)

	// HTTP Server
	handler := gin.Default()
	v1.NewRouter(handler, productUsecase, l)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: ", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Info("app - Run - httpServer.Shutdown: %s", err)
	}
}
