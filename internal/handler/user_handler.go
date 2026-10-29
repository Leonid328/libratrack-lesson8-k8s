package handler

import (
	"net/http"

	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	users repository.UserRepository
}

func NewUserHandler(users repository.UserRepository) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) List(c *gin.Context) {
	list, err := h.users.List(c.Request.Context())
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		writeErr(c, repository.ErrNotFound)
		return
	}
	var in struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		writeErr(c, repository.ErrValidation)
		return
	}
	u, err := h.users.UpdateRole(c.Request.Context(), id, in.Role)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) MyLoans(c *gin.Context) {
	c.JSON(http.StatusOK, []any{})
}
