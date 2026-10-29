package handler

import (
	"net/http"

	"libratrack/internal/model"
	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var cat model.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		writeErr(c, repository.ErrValidation)
		return
	}
	out, err := h.repo.Create(c.Request.Context(), cat)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *CategoryHandler) Get(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	out, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.repo.List(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
