package app

import (
	"log"
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
	"github.com/idoyudha/eshop-product/pkg/redis"
)

func Run(cfg *config.Config) {
	dynamoDB, err := dynamodb.NewDynamoDB(&cfg.DynamoDB)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewRedis(cfg.Redis)

	productUsecase := usecase.New(
		repo.NewProductRepo(dynamoDB),
		repo.NewCategoryDynamoRepo(dynamoDB),
		repo.NewCategoryRedisRepo(redisClient),
	)

	// HTTP Server
	handler := gin.Default()
	v1.NewRouter(handler, productUsecase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: ", s.String())
	case err = <-httpServer.Notify():
		log.Println("app - Run - httpServer.Notify: ", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: ", err)
	}
}
