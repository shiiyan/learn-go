package main

import (
	"context"
	"testing"
)

func TestPersonService_Register(t *testing.T) {
	var gotPerson *Person
	var gotConfirm bool

	mockStore := &PersonStoreMock{
		CreateFunc: func(ctx context.Context, person *Person, confirm bool) error {
			gotPerson = person
			gotConfirm = confirm

			return nil
		},
	}

	s := NewPersonService(mockStore)
	ctx := context.Background()
	p := &Person{ID: "100", Name: "Alice"}

	if err := s.Register(ctx, p); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if len(mockStore.CreateCalls()) != 1 {
		t.Error("expected Create to be called")
	}
	if gotPerson != p {
		t.Errorf("expected %v, got %v", p, gotPerson)
	}
	if !gotConfirm {
		t.Error("expected true, got false")
	}
}
