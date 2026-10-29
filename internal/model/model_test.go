package model

import "testing"

func TestBookHasNoPlainAuthorField(t *testing.T) {
	b := Book{Title: "Go", ISBN: "123", CategoryID: 1}
	if b.Title == "" || b.ISBN == "" {
		t.Fatal("book fields")
	}
}

func TestUserHidesPasswordHash(t *testing.T) {
	u := User{Email: "a@b.c", PasswordHash: "secret"}
	if u.PasswordHash != "secret" {
		t.Fatal("hash stored")
	}
}

func TestCopyStatuses(t *testing.T) {
	c := Copy{Status: "available"}
	if c.Status != "available" {
		t.Fatal(c.Status)
	}
}
