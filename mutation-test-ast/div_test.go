package main

import (
	"errors"
	"testing"
)

// go test -v --overlay overlay.json

func TestDivide(t *testing.T) {
	expected := 1
	res, err := divide(1, 1)

	if err != nil {
		t.Fatalf("expected no error, but go %s", err)
	}

	if res != expected {
		t.Fatalf("expected result %d, but got %d", expected, res)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := divide(1, 0)
	if !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("expected %s, but got %s", ErrDivideByZero, err)
	}
}
