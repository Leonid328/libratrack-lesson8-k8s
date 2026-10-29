package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"libratrack/internal/auth"
	"libratrack/internal/middleware"
	"libratrack/internal/model"
	"libratrack/internal/repository"
	appweb "libratrack/web"

	"github.com/gin-gonic/gin"
)

type fakeBooks struct {
	items []model.Book
}

func (f *fakeBooks) Create(ctx context.Context, in model.BookInput) (model.Book, error) {
	return model.Book{}, nil
}
func (f *fakeBooks) FindByID(ctx context.Context, id uint) (model.Book, error) {
	return model.Book{}, repository.ErrNotFound
}
func (f *fakeBooks) Update(ctx context.Context, id uint, in model.BookInput) (model.Book, error) {
	return model.Book{}, repository.ErrNotFound
}
func (f *fakeBooks) Delete(ctx context.Context, id uint) error { return nil }
func (f *fakeBooks) List(ctx context.Context, categoryID uint, page, limit int) ([]model.Book, error) {
	return f.items, nil
}

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	tmpl, err := appweb.ParseTemplates()
	if err != nil {
		t.Fatal(err)
	}
	h := NewCatalogHandler(&fakeBooks{items: []model.Book{{
		ID: 1, Title: "The Go Programming Language", ISBN: "9780134190440",
		Category: model.Category{Name: "Programming"},
		Authors:  []model.Author{{FullName: "Alan Donovan"}},
	}}}, tmpl)
	r := gin.New()
	r.GET("/catalog", h.Catalog)
	r.GET("/web/books", h.BookRows)
	r.DELETE("/web/books/:id", middleware.AuthMiddleware(), middleware.RequireRole("librarian", "admin"), h.Delete)
	return r
}

func TestCatalogPage(t *testing.T) {
	r := setupRouter(t)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/catalog", nil))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	body := w.Body.String()
	if !contains(body, "LibraTrack") || !contains(body, "<table") || !contains(body, "The Go Programming Language") {
		t.Fatal(body)
	}
}

func TestHTMXBookSearchFragment(t *testing.T) {
	r := setupRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/web/books?q=Go", nil)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	body := w.Body.String()
	if contains(body, "<html") {
		t.Fatal("expected fragment only")
	}
	if !contains(body, "<tr>") {
		t.Fatal(body)
	}
}

func TestDeleteForbiddenForReader(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	tok, err := auth.IssueToken(2, "reader")
	if err != nil {
		t.Fatal(err)
	}
	r := setupRouter(t)
	req := httptest.NewRequest(http.MethodDelete, "/web/books/1", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatal(w.Code)
	}
}

func TestParseWebID(t *testing.T) {
	id, ok := parseID("9")
	if !ok || id != 9 {
		t.Fatal(id, ok)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(sub) > 0 && indexOf(s, sub) >= 0))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
