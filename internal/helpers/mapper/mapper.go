package mapper

import (
	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/entity"
)

func ProductRequestCtxToProductEntity(c *gin.Context) *entity.Product {
	id := c.Param("id")
	sku := c.Param("sku")
	name := c.Param("name")
	imageURL := c.Param("image_url")
	description := c.Param("description")
	price := c.GetFloat64("price")
	quantity := c.GetInt("quantity")
	categoryId := c.GetInt("category_id")

	return &entity.Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		ImageURL:    imageURL,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CategoryID:  categoryId,
	}
}
