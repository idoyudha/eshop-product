package app

import (
	"github.com/idoyudha/eshop-product/config"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/internal/usecase/repo"
	"github.com/idoyudha/eshop-product/pkg/dynamodb"
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
}
