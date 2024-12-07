package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/helpers/mapper"
	"github.com/idoyudha/eshop-product/internal/usecase"
)

type productRoutes struct {
	uc usecase.ProductUseCase
}

func newProductRoutes(handler *gin.RouterGroup, uc usecase.ProductUseCase) {
	r := &productRoutes{uc: uc}

	h := handler.Group("/products")
	{
		h.POST("/", r.createProduct)
		h.GET("/", r.getProducts)
		h.GET("/:id", r.getProductByID)
		h.GET("/category/:id", r.getProductsByCategory)
		h.PUT("/:id", r.updateProduct)
		h.DELETE("/:id", r.deleteProduct)
	}
}

func (r *productRoutes) createProduct(c *gin.Context) {
	product := mapper.ProductRequestCtxToProductEntity(c)
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	err := r.uc.CreateProduct(c.Request.Context(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "success create products")
}

func (r *productRoutes) getProducts(c *gin.Context) {
	products, err := r.uc.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (r *productRoutes) getProductByID(c *gin.Context) {
	product, err := r.uc.GetProductByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (r *productRoutes) getProductsByCategory(c *gin.Context) {
	products, err := r.uc.GetProductsByCategory(c.Request.Context(), c.GetInt("category_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (r *productRoutes) updateProduct(c *gin.Context) {
	product := mapper.ProductRequestCtxToProductEntity(c)
	product.UpdatedAt = time.Now()
	err := r.uc.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success update product")
}

func (r *productRoutes) deleteProduct(c *gin.Context) {
	err := r.uc.DeleteProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success delete product")
}
