package web

import (
	"html/template"
	"net/http"
	"strings"

	"libratrack/internal/model"
	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	books repository.BookRepository
	tmpl  *template.Template
}

func NewCatalogHandler(books repository.BookRepository, tmpl *template.Template) *CatalogHandler {
	return &CatalogHandler{books: books, tmpl: tmpl}
}

func (h *CatalogHandler) Catalog(c *gin.Context) {
	books, err := h.books.List(c.Request.Context(), 0, 1, 100)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = h.tmpl.ExecuteTemplate(c.Writer, "catalog", gin.H{"Books": books})
}

func (h *CatalogHandler) BookRows(c *gin.Context) {
	q := strings.ToLower(strings.TrimSpace(c.Query("q")))
	books, err := h.books.List(c.Request.Context(), 0, 1, 100)
	if err != nil {
		c.String(http.StatusInternalServerError, "error")
		return
	}
	var filtered []model.Book
	for _, b := range books {
		if q == "" || strings.Contains(strings.ToLower(b.Title), q) || strings.Contains(strings.ToLower(b.ISBN), q) {
			filtered = append(filtered, b)
		}
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	for _, b := range filtered {
		_ = h.tmpl.ExecuteTemplate(c.Writer, "book_row", b)
	}
}

func (h *CatalogHandler) Delete(c *gin.Context) {
	id, ok := parseID(c.Param("id"))
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	if err := h.books.Delete(c.Request.Context(), id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusOK)
}

func parseID(s string) (uint, bool) {
	var n uint
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		n = n*10 + uint(ch-'0')
	}
	if n == 0 {
		return 0, false
	}
	return n, true
}
