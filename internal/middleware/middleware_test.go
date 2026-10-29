package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"libratrack/internal/auth"

	"github.com/gin-gonic/gin"
)

func TestRedactPassword(t *testing.T) {
	got := redact([]byte(`{"email":"a@b.c","password":"x"}`))
	if got == "" || contains(got, `"x"`) {
		t.Fatal(got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || (len(s) > 0 && indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func TestRequireRoleForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", "reader"); c.Next() })
	r.GET("/x", RequireRole("admin"), func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Code != http.StatusForbidden {
		t.Fatal(w.Code)
	}
}

func TestAuthMiddlewareMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", AuthMiddleware(), func(c *gin.Context) { c.Status(200) })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatal(w.Code)
	}
}

func TestAuthMiddlewareBearer(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	tok, err := auth.IssueToken(1, "librarian")
	if err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", AuthMiddleware(), RequireRole("librarian", "admin"), func(c *gin.Context) { c.Status(201) })
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatal(w.Body.String(), w.Code)
	}
}
