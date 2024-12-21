package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/config"
	v1Http "github.com/idoyudha/eshop-product/internal/controller/http/v1"
	v1Kafka "github.com/idoyudha/eshop-product/internal/controller/kafka/v1"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/internal/usecase/repo"
	"github.com/idoyudha/eshop-product/pkg/dynamodb"
	"github.com/idoyudha/eshop-product/pkg/httpserver"
	"github.com/idoyudha/eshop-product/pkg/kafka"
	"github.com/idoyudha/eshop-product/pkg/logger"
	"github.com/idoyudha/eshop-product/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Broker)
	if err != nil {
		l.Fatal("app - Run - kafka.NewKafkaProducer: ", err)
	}
	defer kafkaProducer.Close()

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Broker)
	if err != nil {
		l.Fatal("app - Run - kafka.NewKafkaConsumer: ", err)
	}
	defer kafkaConsumer.Close()

	dynamoDB, err := dynamodb.NewDynamoDB(&cfg.AWS)
	if err != nil {
		l.Fatal("app - Run - dynamodb.NewDynamoDB: ", err)
	}

	redisClient, err := redis.NewRedis(cfg.Redis)
	if err != nil {
		l.Fatal("app - Run - redis.NewRedis: ", err)
	}

	productUseCase := usecase.NewProductUseCase(
		repo.NewProductRepo(dynamoDB),
		kafkaProducer,
	)

	categoryUseCase := usecase.NewCategoryUseCase(
		repo.NewCategoryRedisRepo(redisClient),
		repo.NewCategoryDynamoRepo(dynamoDB),
	)

	// HTTP Server
	handler := gin.Default()
	v1Http.HTTPNewRouter(handler, productUseCase, categoryUseCase, l)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Kafka Consumer
	kafkaErrChan := make(chan error, 1)
	go func() {
		if err := v1Kafka.KafkaNewRouter(productUseCase, l, kafkaConsumer); err != nil {
			kafkaErrChan <- err
		}
	}()

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
