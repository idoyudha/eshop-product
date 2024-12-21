package usecase

import (
	"context"

	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/redis/go-redis/v9"
)

type CategoryUseCase struct {
	categoryRepoDynamo CategoryDynamoRepo
	categoryRepoRedis  CategoryRedisRepo
}

func NewCategoryUseCase(
	categoryRepoRedis CategoryRedisRepo,
	categoryRepoDynamo CategoryDynamoRepo,
) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepoRedis:  categoryRepoRedis,
		categoryRepoDynamo: categoryRepoDynamo,
	}
}

func (u *CategoryUseCase) GetCategories(ctx context.Context) (*[]entity.Category, error) {
	// get from redis first
	categories, err := u.categoryRepoRedis.GetAll(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if categories != nil {
		return categories, nil
	}

	// if not found, get from dynamo
	categories, err = u.categoryRepoDynamo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	// set to redis
	err = u.categoryRepoRedis.SaveAll(ctx, categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (u *CategoryUseCase) CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	err := category.GenerateCategoryID()
	if err != nil {
		return nil, err
	}

	// create new in dynamodb
	err = u.categoryRepoDynamo.Save(ctx, category)
	if err != nil {
		return nil, err
	}

	// set new in redis
	err = u.categoryRepoRedis.Add(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *CategoryUseCase) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// update in dynamodb
	err := u.categoryRepoDynamo.Update(ctx, category)
	if err != nil {
		return err
	}

	// update in redis
	err = u.categoryRepoRedis.Update(ctx, category.ID, category.Name)
	if err != nil {
		return err
	}

	return nil
}

func (u *CategoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	// delete in dynamodb
	err := u.categoryRepoDynamo.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = u.categoryRepoRedis.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
