package v1

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

func HTTPNewRouter(
	handler *gin.Engine,
	ucp usecase.Product,
	ucg usecase.Category,
	l logger.Interface,
) {
	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

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
