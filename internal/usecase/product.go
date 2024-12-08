package usecase

import (
	"context"
	"log"

	"github.com/idoyudha/eshop-product/internal/entity"
)

type ProductUseCase struct {
	productRepoDynamo  ProductDynamoRepo
	categoryRepoDynamo CategoryDynamoRepo
	categoryRepoRedis  CategoryRedisRepo
}

func NewProductUseCase(productRepoDynamo ProductDynamoRepo) *ProductUseCase {
	return &ProductUseCase{
		productRepoDynamo: productRepoDynamo,
	}
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, product *entity.Product) error {
	log.Println("CreateProduct", product)
	return u.productRepoDynamo.Save(ctx, product)
}

func (u *ProductUseCase) GetProducts(ctx context.Context) (*[]entity.Product, error) {
	return u.productRepoDynamo.GetProducts(ctx)
}

func (u *ProductUseCase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return u.productRepoDynamo.GetProductByID(ctx, id)
}

func (u *ProductUseCase) GetProductsByCategory(ctx context.Context, categoryID int) ([]entity.Product, error) {
	return u.productRepoDynamo.GetProductsByCategory(ctx, categoryID)
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return u.productRepoDynamo.Update(ctx, product)
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, productID string, categoryID int) error {
	return u.productRepoDynamo.Delete(ctx, productID, categoryID)
}
