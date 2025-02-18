package v1

import (
	"errors"
	"mime/multipart"
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
		h.POST("", r.createProduct)
		h.GET("", r.getProducts)
		h.GET("/:id", r.getProductByID)
		h.GET("/category/:id", r.getProductsByCategory)
		h.POST("/categories", r.getProductsByCategories)
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
	if err := c.ShouldBind(&request); err != nil {
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

type getProductsRequest struct {
	CategoryIDs []string `json:"category_ids" binding:"required"`
}

func (r *productRoutes) getProductsByCategories(c *gin.Context) {
	var request getProductsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	productEntities, err := r.uc.GetProductsByCategories(c.Request.Context(), request.CategoryIDs)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
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
	CategoryID  string                `form:"category_id" binding:"required"`
}

func (u *updateProductRequest) validate() error {
	if u.Image.Size > 1024*1024 {
		return errors.New("image size must be less than 1MB")
	}

	contentType := u.Image.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return errors.New("image must be in JPEG or PNG format")
	}

	return nil
}

type updateProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
}

func (r *productRoutes) updateProduct(c *gin.Context) {
	var request updateProductRequest
	if err := c.ShouldBind(&request); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	if err := request.validate(); err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	productEntity := updateProductRequestToProductEntity(request, c.Param("id"))

	err := r.uc.UpdateProduct(c.Request.Context(), &productEntity, request.Image)
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - updateProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	product := productEntityToUpdateProductResponse(productEntity)

	c.JSON(http.StatusOK, newUpdateSuccess(product))
}

func (r *productRoutes) deleteProduct(c *gin.Context) {
	err := r.uc.DeleteProduct(c.Request.Context(), c.Param("product_id"), c.Param("category_id"))
	if err != nil {
		r.l.Error(err, "http - v1 - productRoutes - deleteProduct")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newDeleteSuccess())
}
