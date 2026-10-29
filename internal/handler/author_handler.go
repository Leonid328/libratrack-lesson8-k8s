package handler

import (
	"net/http"

	"libratrack/internal/model"
	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	repo repository.AuthorRepository
}

func NewAuthorHandler(repo repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{repo: repo}
}

func (h *AuthorHandler) Create(c *gin.Context) {
	var a model.Author
	if err := c.ShouldBindJSON(&a); err != nil {
		writeErr(c, repository.ErrValidation)
		return
	}
	out, err := h.repo.Create(c.Request.Context(), a)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *AuthorHandler) Get(c *gin.Context) {
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

func (h *AuthorHandler) List(c *gin.Context) {
	list, err := h.repo.List(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
