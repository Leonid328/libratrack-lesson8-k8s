package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"libratrack/internal/auth"

	"github.com/gin-gonic/gin"
)

func TestLoginRejectsBadJSON(t *testing.T) {
	h := NewAuthHandler(nil)
	r := gin.New()
	r.POST("/login", h.Login)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatal(w.Code)
	}
}

func TestRegisterValidation(t *testing.T) {
	h := NewAuthHandler(nil)
	r := gin.New()
	r.POST("/register", h.Register)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 422 {
		t.Fatal(w.Code, w.Body.String())
	}
}

func TestIssuedJWTHasExp(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	tok, err := auth.IssueToken(1, "reader")
	if err != nil {
		t.Fatal(err)
	}
	c, err := auth.ParseToken(tok)
	if err != nil || c.ExpiresAt == nil {
		t.Fatal(err, c)
	}
	_ = json.RawMessage{}
}
