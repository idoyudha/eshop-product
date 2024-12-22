package v1

import (
	"mime/multipart"
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

type createProductRequest struct {
	Name        string                `form:"name" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Price       float64               `form:"price" binding:"required"`
	Quantity    int                   `form:"quantity" binding:"required"`
	CategoryID  string                `form:"category_id" binding:"required"`
}

type createProductResponse struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  string  `json:"category_id"`
}

func (r *productRoutes) createProduct(c *gin.Context) {
	var request createProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusBadRequest, newInternalServerError(err.Error()))
		return
	}

	productEntity := createProductRequestToProductEntity(request)

	product, err := r.uc.CreateProduct(c.Request.Context(), &productEntity, request.Image)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - createProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, newCreateSuccess(product))
}

type getProductResponse struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  string  `json:"category_id"`
}

func (r *productRoutes) getProducts(c *gin.Context) {
	products, err := r.uc.GetProducts(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProducts")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	productsResponse := productEntitiesToGetProductResponse(*products)

	c.JSON(http.StatusOK, newGetSuccess(productsResponse))
}

func (r *productRoutes) getProductByID(c *gin.Context) {
	product, err := r.uc.GetProductByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductByID")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	productsResponse := productEntityToGetProductResponse(*product)

	c.JSON(http.StatusOK, newGetSuccess(productsResponse))
}

func (r *productRoutes) getProductsByCategory(c *gin.Context) {
	productEntities, err := r.uc.GetProductsByCategory(c.Request.Context(), c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - getProductsByCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	products := productEntitiesToGetProductResponse(productEntities)

	c.JSON(http.StatusOK, newGetSuccess(products))
}

type updateProductRequest struct {
	Name        string                `form:"name" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Price       float64               `form:"price" binding:"required"`
	Quantity    int                   `form:"quantity" binding:"required"`
	CategoryID  string                `form:"category_id" binding:"required"`
}

type updateProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  string  `json:"category_id"`
}

func (r *productRoutes) updateProduct(c *gin.Context) {
	var request updateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	productEntity := updateProductRequestToProductEntity(request, c.Param("id"))

	err := r.uc.UpdateProduct(c.Request.Context(), &productEntity)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	product := productEntityToUpdateProductResponse(productEntity)

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
