package handler

import (
	"net/http"

	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type CopyHandler struct {
	repo repository.CopyRepository
}

func NewCopyHandler(repo repository.CopyRepository) *CopyHandler {
	return &CopyHandler{repo: repo}
}

type copyInput struct {
	InventoryNo string `json:"inventory_no"`
	Status      string `json:"status"`
}

func (h *CopyHandler) Create(c *gin.Context) {
	bookID, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	var in copyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, repository.ErrValidation)
		return
	}
	cp, err := h.repo.Create(c.Request.Context(), bookID, in.InventoryNo, in.Status)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, cp)
}

func (h *CopyHandler) List(c *gin.Context) {
	bookID, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	list, err := h.repo.ListByBook(c.Request.Context(), bookID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CopyHandler) Availability(c *gin.Context) {
	bookID, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	av, err := h.repo.Availability(c.Request.Context(), bookID)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, av)
}
