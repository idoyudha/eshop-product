package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/usecase"
)

func NewRouter(handler *gin.Engine, uc usecase.Product) {
	// options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// health check
	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	h := handler.Group("/v1")
	{
		newProductRoutes(h, uc)
	}
}
