package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-product/internal/entity"
	"github.com/idoyudha/eshop-product/internal/usecase"
	"github.com/idoyudha/eshop-product/pkg/logger"
)

type categoryRoutes struct {
	uc usecase.Category
	l  logger.Interface
}

func newCategoryRoutes(handler *gin.RouterGroup, uc usecase.Category, l logger.Interface) {
	r := &categoryRoutes{uc: uc, l: l}

	h := handler.Group("/categories")
	{
		h.POST("/", r.createCategory)
		h.GET("/", r.getCategories)
		h.PUT("/:id", r.updateCategory)
		h.DELETE("/:id", r.deleteCategory)
	}
}

type createCategoryRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parent_id" binding:"required"`
}

type createCategoryResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

func (r *categoryRoutes) createCategory(c *gin.Context) {
	var request createCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - createCategory")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	categoryEntity := createCategoryRequestToCategoryEntity(request)

	category, err := r.uc.CreateCategory(c.Request.Context(), &categoryEntity)
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - createCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	categoryResponse := categoryEntityToCreateCategoryResponse(*category)

	c.JSON(http.StatusCreated, newCreateSuccess(categoryResponse))
}

type getCategoryResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

func (r *categoryRoutes) getCategories(c *gin.Context) {
	categories, err := r.uc.GetCategories(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - getCategories")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	categoriesResponse := categoryEntitiesToGetCategoryResponse(*categories)

	c.JSON(http.StatusOK, newGetSuccess(categoriesResponse))
}

type updateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (r *categoryRoutes) updateCategory(c *gin.Context) {
	var request updateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - updateCategory")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	category := entity.Category{
		ID:   c.Param("id"),
		Name: request.Name,
	}

	err := r.uc.UpdateCategory(c.Request.Context(), &category)
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - updateCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	categoryResponse := categoryEntityToUpdateCategoryResponse(category)

	c.JSON(http.StatusOK, newUpdateSuccess(categoryResponse))
}

func (r *categoryRoutes) deleteCategory(c *gin.Context) {
	err := r.uc.DeleteCategory(c.Request.Context(), c.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - deleteCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newDeleteSuccess())
}
