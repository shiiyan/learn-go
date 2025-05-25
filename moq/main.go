package main

import (
	"context"
	"errors"
)

type Person struct {
	ID      string
	Name    string
	Company string
	Website string
}

type PersonStore interface {
	Get(ctx context.Context, id string) (*Person, error)
	Create(ctx context.Context, person *Person, confirm bool) error
}

type PersonService struct {
	store PersonStore
}

func NewPersonService(store PersonStore) *PersonService {
	return &PersonService{store: store}
}

func (s *PersonService) Register(ctx context.Context, p *Person) error {
	if p == nil {
		return errors.New("person is nil")
	}

	return s.store.Create(ctx, p, true)
}

func (s *PersonService) GetName(ctx context.Context, id string) (string, error) {
	p, err := s.store.Get(ctx, id)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", errors.New("person not found")
	}

	return p.Name, nil
}
