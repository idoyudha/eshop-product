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

type CreateCategoryRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

func (r *categoryRoutes) createCategory(c *gin.Context) {
	var request CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - createCategory")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	category, err := CreateCategoryRequestToCategoryEntity(request)
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - createCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	err = r.uc.CreateCategory(c.Request.Context(), &category)
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - createCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, newCreateSuccess(category))
}

func (r *categoryRoutes) getCategories(c *gin.Context) {
	categories, err := r.uc.GetCategories(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - getCategories")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newGetSuccess(categories))
}

func (r *categoryRoutes) updateCategory(c *gin.Context) {
	var request CreateCategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - updateCategory")
		c.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	category := entity.Category{
		ID:       c.Param("id"),
		Name:     request.Name,
		ParentID: request.ParentID,
	}

	err := r.uc.UpdateCategory(c.Request.Context(), &category)
	if err != nil {
		r.l.Error(err, "http - v1 - categoryRoutes - updateCategory")
		c.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, newUpdateSuccess(category))
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
