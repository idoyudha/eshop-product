package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/idoyudha/eshop-product/internal/entity"
)

type ProductDynamoRepo struct {
	db *dynamodb.Client
}

func NewProductRepo(db *dynamodb.Client) *ProductDynamoRepo {
	return &ProductDynamoRepo{
		db: db,
	}
}

func (r *ProductDynamoRepo) Save(ctx context.Context, product *entity.Product) error {
	// TODO: implement save product
	return nil
}

func (r *ProductDynamoRepo) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	// TODO: implement scan all of products
	return nil, nil
}

func (r *ProductDynamoRepo) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	// TODO: implement to get product by id
	return nil, nil
}

func (r *ProductDynamoRepo) GetProductsByCategory(ctx context.Context, categoryID int) ([]entity.Product, error) {
	// TODO: implement to query products by category_id
	return nil, nil
}

func (r *ProductDynamoRepo) Update(ctx context.Context, product *entity.Product) error {
	// TODO: implement update product
	return nil
}

func (r *ProductDynamoRepo) Delete(ctx context.Context, id string) error {
	// TODO: implement delete product
	return nil
}
