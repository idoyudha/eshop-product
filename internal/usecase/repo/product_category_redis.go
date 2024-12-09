package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/idoyudha/eshop-product/pkg/redis"
)

const (
	categoriesKey = "categories"
)

type CategoryRedisRepo struct {
	*redis.RedisClient
}

func NewCategoryRedisRepo(redis *redis.RedisClient) *CategoryRedisRepo {
	return &CategoryRedisRepo{
		redis,
	}
}

func (r *CategoryRedisRepo) Save(ctx context.Context, categories *[]entity.Category) error {
	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return fmt.Errorf("failed to marshal categories: %w", err)
	}

	err = r.Client.Set(ctx, categoriesKey, categoriesJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save categories to redis: %w", err)
	}

	return nil
}

func (r *CategoryRedisRepo) GetCategories(ctx context.Context) (*[]entity.Category, error) {
	categoriesJSON, err := r.Client.Get(ctx, categoriesKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from redis: %w", err)
	}

	var categories []entity.Category
	if err := json.Unmarshal([]byte(categoriesJSON), &categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
	}

	return &categories, nil
}

func (r *CategoryRedisRepo) Delete(ctx context.Context) error {
	if err := r.Client.Del(ctx, categoriesKey).Err(); err != nil {
		return fmt.Errorf("failed to delete categories from redis: %w", err)
	}

	return nil
}
