package handler

import (
	"net/http"
	"strconv"

	"libratrack/internal/model"
	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

func (h *BookHandler) Create(c *gin.Context) {
	var in model.BookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, repository.ErrValidation)
		return
	}
	book, err := h.repo.Create(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) Get(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	book, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) List(c *gin.Context) {
	var cat uint
	if v := c.Query("category_id"); v != "" {
		n, _ := strconv.ParseUint(v, 10, 64)
		cat = uint(n)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	books, err := h.repo.List(c.Request.Context(), cat, page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) Update(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	var in model.BookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, repository.ErrValidation)
		return
	}
	book, err := h.repo.Update(c.Request.Context(), id, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) Delete(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		writeErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(s string) (uint, bool) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil || n == 0 {
		return 0, false
	}
	return uint(n), true
}
