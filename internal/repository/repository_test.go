package repository

import (
	"testing"

	"libratrack/internal/model"
)

func TestErrNotFoundDeclared(t *testing.T) {
	if ErrNotFound.Error() != "not found" {
		t.Fatal(ErrNotFound)
	}
}

func TestErrConflict(t *testing.T) {
	if ErrConflict.Error() != "conflict" {
		t.Fatal(ErrConflict)
	}
}

func TestUniqueIDs(t *testing.T) {
	got := uniqueIDs([]uint{1, 1, 2})
	if len(got) != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestValidateBookInputEmpty(t *testing.T) {
	if err := validateBookInput(model.BookInput{}); err == nil {
		t.Fatal("expected validation")
	}
}

func TestValidateBookInputOK(t *testing.T) {
	err := validateBookInput(model.BookInput{Title: "T", ISBN: "1", CategoryID: 1, AuthorIDs: []uint{1}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(ErrNotFound) {
		t.Fatal("expected true")
	}
}
