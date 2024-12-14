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
		Delete(context.Context, string, int) error
	}

	CategoryDynamoRepo interface {
		Save(context.Context, *entity.Category) error
		GetCategories(context.Context) (*[]entity.Category, error)
		Update(context.Context, *entity.Category) error
		Delete(context.Context, string) error
	}

	CategoryRedisRepo interface {
		SaveAll(context.Context, *[]entity.Category) error
		GetAll(context.Context) (*[]entity.Category, error)
		GetByParentID(context.Context, string) ([]entity.Category, error)
		Add(context.Context, *entity.Category) error
		UpdateName(context.Context, string, string) error
		Delete(context.Context, string) error
	}

	Product interface {
		CreateProduct(context.Context, *entity.Product) error
		GetProducts(context.Context) (*[]entity.Product, error)
		GetProductByID(context.Context, string) (*entity.Product, error)
		GetProductsByCategory(context.Context, int) ([]entity.Product, error)
		UpdateProduct(context.Context, *entity.Product) error
		DeleteProduct(context.Context, string, int) error
	}

	Category interface {
		CreateCategory(context.Context, *entity.Category) error
		GetCategories(context.Context) (*[]entity.Category, error)
		UpdateCategory(context.Context, *entity.Category) error
		DeleteCategory(context.Context, string) error
	}
)
