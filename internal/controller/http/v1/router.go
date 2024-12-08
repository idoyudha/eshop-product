package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

func NewRouter(
	handler *gin.Engine,
	ucp usecase.Product,
	ucg usecase.Category,
	l logger.Interface,
) {
	// options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// health check
	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	h := handler.Group("/v1")
	{
		newProductRoutes(h, ucp, l)
		newCategoryRoutes(h, ucg, l)
	}
}
