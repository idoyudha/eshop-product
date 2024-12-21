package v1

import (
	"time"

	"github.com/idoyudha/eshop-product/internal/entity"
)

func createProductRequestToProductEntity(request createProductRequest) entity.Product {
	return entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Quantity:    request.Quantity,
		CategoryID:  request.CategoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func productEntityToProductResponse(product entity.Product) createProductResponse {
	return createProductResponse{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		ImageURL:    product.ImageURL,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
	}
}

func updateProductRequestToProductEntity(request updateProductRequest, id string) entity.Product {
	return entity.Product{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Quantity:    request.Quantity,
		CategoryID:  request.CategoryID,
		UpdatedAt:   time.Now(),
	}
}

func productEntityToUpdateProductResponse(product entity.Product) updateProductResponse {
	return updateProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		ImageURL:    product.ImageURL,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
	}
}

func productEntitiesToGetProductResponse(product []entity.Product) []getProductResponse {
	var response []getProductResponse
	for _, p := range product {
		response = append(response, getProductResponse{
			ID:          p.ID,
			SKU:         p.SKU,
			Name:        p.Name,
			ImageURL:    p.ImageURL,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			CategoryID:  p.CategoryID,
		})
	}
	return response
}

func productEntityToGetProductResponse(product entity.Product) getProductResponse {
	return getProductResponse{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		ImageURL:    product.ImageURL,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
	}
}

func createCategoryRequestToCategoryEntity(request createCategoryRequest) entity.Category {
	return entity.Category{
		Name:      request.Name,
		ParentID:  request.ParentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func categoryEntityToCreateCategoryResponse(category entity.Category) createCategoryResponse {
	return createCategoryResponse{
		ID:       category.ID,
		Name:     category.Name,
		ParentID: category.ParentID,
	}
}

func categoryEntitiesToGetCategoryResponse(categories []entity.Category) []getCategoryResponse {
	var response []getCategoryResponse
	for _, c := range categories {
		response = append(response, getCategoryResponse{
			ID:       c.ID,
			Name:     c.Name,
			ParentID: c.ParentID,
		})
	}
	return response
}

func categoryEntityToUpdateCategoryResponse(category entity.Category) updateCategoryResponse {
	return updateCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}
