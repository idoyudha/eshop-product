package usecase

import (
	"context"

	"github.com/idoyudha/eshop-product/internal/entity"
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
	categories, err := u.categoryRepoRedis.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	if len(*categories) > 0 {
		return categories, nil
	}

	// if not found, get from dynamo
	categories, err = u.categoryRepoDynamo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	// set to redis
	err = u.categoryRepoRedis.Save(ctx, categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (u *CategoryUseCase) CreateCategory(ctx context.Context, category *entity.Category) error {
	// create new in dynamodb
	err := u.categoryRepoDynamo.Save(ctx, category)
	if err != nil {
		return err
	}

	// get all in dynamodb
	categories, err := u.categoryRepoDynamo.GetCategories(ctx)
	if err != nil {
		return err
	}

	// set new in redis, replace the last one
	err = u.categoryRepoRedis.Save(ctx, categories)
	if err != nil {
		return err
	}

	return nil
}

func (u *CategoryUseCase) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// update in dynamodb
	err := u.categoryRepoDynamo.Update(ctx, category)
	if err != nil {
		return err
	}

	// get all in dynamodb
	categories, err := u.categoryRepoDynamo.GetCategories(ctx)
	if err != nil {
		return err
	}

	// set new in redis, replace the last one
	err = u.categoryRepoRedis.Save(ctx, categories)
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

	// get all in dynamodb
	categories, err := u.categoryRepoDynamo.GetCategories(ctx)
	if err != nil {
		return err
	}

	// set new in redis, replace the last one
	if len(*categories) == 0 {
		// delete the value in redis
		err = u.categoryRepoRedis.Delete(ctx)
		if err != nil {
			return err
		}
	} else {
		err = u.categoryRepoRedis.Save(ctx, categories)
		if err != nil {
			return err
		}
	}

	return nil
}
