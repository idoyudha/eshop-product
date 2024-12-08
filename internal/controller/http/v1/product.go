package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

type productRoutes struct {
	uc usecase.Product
	l  logger.Interface
}

func newProductRoutes(handler *gin.RouterGroup, uc usecase.Product, l logger.Interface) {
	r := &productRoutes{uc: uc, l: l}

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

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	ImageURL    string  `json:"image_url" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
}

func (r *productRoutes) createProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusBadRequest, response{Error: err.Error()})
		return
	}

	product, err := CreateProductRequestToProductEntity(request)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	err = r.uc.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "success create products")
}

func (r *productRoutes) getProducts(c *gin.Context) {
	products, err := r.uc.GetProducts(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProducts")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (r *productRoutes) getProductByID(c *gin.Context) {
	product, err := r.uc.GetProductByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductByID")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (r *productRoutes) getProductsByCategory(c *gin.Context) {
	products, err := r.uc.GetProductsByCategory(c.Request.Context(), c.GetInt("category_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductsByCategory")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	ImageURL    string  `json:"image_url" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
}

func (r *productRoutes) updateProduct(c *gin.Context) {
	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusBadRequest, response{Error: err.Error()})
		return
	}

	product := UpdateProductRequestToProductEntity(request)
	err := r.uc.UpdateProduct(c.Request.Context(), &product)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success update product")
}

func (r *productRoutes) deleteProduct(c *gin.Context) {
	err := r.uc.DeleteProduct(c.Request.Context(), c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - deleteProduct")
		c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success delete product")
}
