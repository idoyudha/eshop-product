package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/idoyudha/eshop-product/internal/utils"
)

func CreateProductRequestToProductEntity(request CreateProductRequest) (entity.Product, error) {
	productId, err := uuid.NewV7()
	if err != nil {
		return entity.Product{}, err
	}
	return entity.Product{
		ID:          productId.String(),
		SKU:         utils.GenerateSKU(),
		Name:        request.Name,
		ImageURL:    request.ImageURL,
		Description: request.Description,
		Price:       request.Price,
		Quantity:    request.Quantity,
		CategoryID:  request.CategoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func UpdateProductRequestToProductEntity(request UpdateProductRequest) entity.Product {
	return entity.Product{
		Name:        request.Name,
		ImageURL:    request.ImageURL,
		Description: request.Description,
		Price:       request.Price,
		Quantity:    request.Quantity,
		CategoryID:  request.CategoryID,
		UpdatedAt:   time.Now(),
	}
}
