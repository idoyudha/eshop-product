package repo

import (
	"context"

	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/idoyudha/eshop-product/pkg/dynamodb"
)

type CategoryDynamoRepo struct {
	*dynamodb.DynamoDB
}

func NewCategoryDynamoRepo(db *dynamodb.DynamoDB) *CategoryDynamoRepo {
	return &CategoryDynamoRepo{
		db,
	}
}

func (r *CategoryDynamoRepo) Save(ctx context.Context, category *entity.Category) error {
	// TODO: implement save category
	return nil
}

func (r *CategoryDynamoRepo) GetCategories(ctx context.Context) (*[]entity.Category, error) {
	// TODO: implement scan all of categories
	return nil, nil
}

func (r *CategoryDynamoRepo) Update(ctx context.Context, category *entity.Category) error {
	// TODO: implement update category
	return nil
}

func (r *CategoryDynamoRepo) Delete(ctx context.Context, id string) error {
	// TODO: implement delete category
	return nil
}
