package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRoleUnauthorizedWithoutRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", RequireRole("admin"), func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatal(w.Code)
	}
}

func TestRequireRoleAllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", "admin"); c.Next() })
	r.POST("/books", RequireRole("librarian", "admin"), func(c *gin.Context) { c.Status(http.StatusCreated) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/books", nil))
	if w.Code != http.StatusCreated {
		t.Fatal(w.Code)
	}
}
