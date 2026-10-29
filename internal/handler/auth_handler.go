package handler

import (
	"net/http"
	"strings"

	"libratrack/internal/auth"
	"libratrack/internal/model"
	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	users repository.UserRepository
}

func NewAuthHandler(users repository.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var in registerRequest
	if err := c.ShouldBindJSON(&in); err != nil || strings.TrimSpace(in.Email) == "" || in.Password == "" {
		writeErr(c, repository.ErrValidation)
		return
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		writeErr(c, err)
		return
	}
	u, err := h.users.Create(c.Request.Context(), model.User{
		Email:        in.Email,
		PasswordHash: hash,
		Role:         "reader",
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in loginRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "unauthorized", Message: "invalid credentials"}})
		return
	}
	u, err := h.users.FindByEmail(c.Request.Context(), in.Email)
	if err != nil || auth.CheckPassword(u.PasswordHash, in.Password) != nil {
		c.JSON(http.StatusUnauthorized, errorBody{Error: errorDetail{Code: "unauthorized", Message: "invalid credentials"}})
		return
	}
	token, err := auth.IssueToken(u.ID, u.Role)
	if err != nil {
		writeErr(c, err)
		return
	}
	c.SetCookie("libratrack_token", token, 86400, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token, "role": u.Role, "user_id": u.ID})
}
