package usecase

import (
	"context"

	"github.com/idoyudha/eshop-product/internal/entity"
)

type (
	ProductDynamoRepo interface {
		Save(context.Context, *entity.Product) error
		GetProducts(context.Context) (*[]entity.Product, error)
		GetProductByID(context.Context, string) (*entity.Product, error)
		GetProductsByCategory(context.Context, int) ([]entity.Product, error)
		Update(context.Context, *entity.Product) error
		Delete(context.Context, string) error
	}

	CategoryDynamoRepo interface {
		Save(context.Context, *entity.Category) error
		GetCategories(context.Context) (*[]entity.Category, error)
		Update(context.Context, *entity.Category) error
		Delete(context.Context, string) error
	}

	CategoryRedisRepo interface {
		Save(context.Context, *entity.Category) error
		GetCategories(context.Context) (*[]entity.Category, error)
		Update(context.Context, *entity.Category) error
		Delete(context.Context, string) error
	}

	Product interface {
		CreateProduct(context.Context, *entity.Product) error
		GetProducts(context.Context) (*[]entity.Product, error)
		GetProductByID(context.Context, string) (*entity.Product, error)
		GetProductsByCategory(context.Context, int) ([]entity.Product, error)
		UpdateProduct(context.Context, *entity.Product) error
		DeleteProduct(context.Context, string) error
	}
)
