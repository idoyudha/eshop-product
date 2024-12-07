package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/idoyudha/eshop-product/internal/entity"
)

type CategoryDynamoRepo struct {
	db *dynamodb.Client
}

func NewCategoryDynamoRepo(db *dynamodb.Client) *CategoryDynamoRepo {
	return &CategoryDynamoRepo{
		db: db,
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

func (r *CategoryDynamoRepo) Delete(ctx context.Context, category *entity.Category) error {
	// TODO: implement update category
	return nil
}

func (r *CategoryDynamoRepo) Update(ctx context.Context, id string) error {
	// TODO: implement delete category
	return nil
}
