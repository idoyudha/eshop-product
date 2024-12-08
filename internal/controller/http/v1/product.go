package v1

import (
	"net/http"
	"strconv"

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
		h.DELETE("/:product_id/category/:category_id", r.deleteProduct)
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
		c.JSON(http.StatusBadRequest, newInternalServerError(err.Error()))
		return
	}

	product, err := CreateProductRequestToProductEntity(request)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	err = r.uc.CreateProduct(c.Request.Context(), &product)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, newCreateSuccess(product))
}

func (r *productRoutes) getProducts(c *gin.Context) {
	products, err := r.uc.GetProducts(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProducts")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newGetSuccess(products))
}

func (r *productRoutes) getProductByID(c *gin.Context) {
	product, err := r.uc.GetProductByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductByID")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newGetSuccess(product))
}

func (r *productRoutes) getProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductsByCategory")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}
	products, err := r.uc.GetProductsByCategory(c.Request.Context(), categoryID)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductsByCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newGetSuccess(products))
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
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	product := UpdateProductRequestToProductEntity(request, c.Param("id"))

	err := r.uc.UpdateProduct(c.Request.Context(), &product)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newUpdateSuccess(product))
}

func (r *productRoutes) deleteProduct(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - deleteProduct")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}
	err = r.uc.DeleteProduct(c.Request.Context(), c.Param("product_id"), categoryID)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - deleteProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newDeleteSuccess())
}
