package auth

import (
	"os"
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	h, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}
	if CheckPassword(h, "secret") != nil {
		t.Fatal("check")
	}
	if CheckPassword(h, "nope") == nil {
		t.Fatal("expected mismatch")
	}
}

func TestIssueAndParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	tok, err := IssueToken(7, "admin")
	if err != nil {
		t.Fatal(err)
	}
	c, err := ParseToken(tok)
	if err != nil {
		t.Fatal(err)
	}
	if c.UserID != 7 || c.Role != "admin" {
		t.Fatalf("%+v", c)
	}
}

func TestParseTokenRejectsGarbage(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	if _, err := ParseToken("not-a-token"); err == nil {
		t.Fatal("expected error")
	}
}
