package handler

import (
	"errors"
	"net/http"

	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, errorBody{Error: errorDetail{Code: "not_found", Message: "not found"}})
	case errors.Is(err, repository.ErrConflict):
		c.JSON(http.StatusConflict, errorBody{Error: errorDetail{Code: "conflict", Message: "conflict"}})
	case errors.Is(err, repository.ErrInvalidRef):
		c.JSON(http.StatusUnprocessableEntity, errorBody{Error: errorDetail{Code: "invalid_reference", Message: "invalid reference"}})
	case errors.Is(err, repository.ErrValidation):
		c.JSON(http.StatusUnprocessableEntity, errorBody{Error: errorDetail{Code: "validation_error", Message: "validation error"}})
	default:
		c.JSON(http.StatusInternalServerError, errorBody{Error: errorDetail{Code: "internal", Message: "internal error"}})
	}
}
