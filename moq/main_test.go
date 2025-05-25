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
		t.Error("expected Create to be called once")
	}
	if gotPerson != p {
		t.Errorf("expected %v, got %v", p, gotPerson)
	}
	if !gotConfirm {
		t.Error("expected true, got false")
	}
}

func TestPersonService_GetName(t *testing.T) {
	expected := &Person{ID: "101", Name: "Bob"}

	mockStore := &PersonStoreMock{
		GetFunc: func(ctx context.Context, id string) (*Person, error) {
			if id != expected.ID {
				t.Errorf("expected id=%q, got %q", expected.ID, id)
			}

			return expected, nil
		},
	}

	s := NewPersonService(mockStore)
	ctx := context.Background()

	name, err := s.GetName(ctx, "101")
	if err != nil {
		t.Fatalf("GetName failed: %v", err)
	}

	if len(mockStore.GetCalls()) != 1 {
		t.Error("expected GetName to be called once")
	}
	if name != expected.Name {
		t.Errorf("expected %q, got %q", expected.Name, name)
	}
}
