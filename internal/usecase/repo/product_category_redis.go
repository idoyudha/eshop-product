package repo

import (
	"context"

	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/redis/go-redis/v9"
)

type CategoryRedisRepo struct {
	redis *redis.Client
}

func NewCategoryRedisRepo(redis *redis.Client) *CategoryRedisRepo {
	return &CategoryRedisRepo{
		redis: redis,
	}
}

func (r *CategoryRedisRepo) Save(ctx context.Context, category *entity.Category) error {
	// TODO: implement save category
	return nil
}

func (r *CategoryRedisRepo) GetCategories(ctx context.Context) (*[]entity.Category, error) {
	// TODO: implement scan all of categories
	return nil, nil
}

func (r *CategoryRedisRepo) Update(ctx context.Context, category *entity.Category) error {
	// TODO: implement update category
	return nil
}

func (r *CategoryRedisRepo) Delete(ctx context.Context, id string) error {
	// TODO: implement delete category
	return nil
}
