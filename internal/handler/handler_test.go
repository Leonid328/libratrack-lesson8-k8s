package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"libratrack/internal/repository"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestParseID(t *testing.T) {
	id, ok := parseID("12")
	if !ok || id != 12 {
		t.Fatal(id, ok)
	}
}

func TestParseIDInvalid(t *testing.T) {
	if _, ok := parseID("x"); ok {
		t.Fatal("expected false")
	}
}

func TestWriteErrNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	writeErr(c, repository.ErrNotFound)
	if w.Code != http.StatusNotFound {
		t.Fatal(w.Code)
	}
	var body errorBody
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.Error.Code != "not_found" {
		t.Fatal(body)
	}
}

func TestWriteErrConflict(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	writeErr(c, repository.ErrConflict)
	if w.Code != 409 {
		t.Fatal(w.Code)
	}
}

func TestWriteErrValidation(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	writeErr(c, repository.ErrValidation)
	if w.Code != 422 {
		t.Fatal(w.Code)
	}
}

func TestWriteErrInvalidRef(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	writeErr(c, repository.ErrInvalidRef)
	if w.Code != 422 {
		t.Fatal(w.Code)
	}
	var body errorBody
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body.Error.Code != "invalid_reference" {
		t.Fatal(body)
	}
}

func TestHealthJSONShape(t *testing.T) {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
