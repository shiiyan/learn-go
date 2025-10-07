package main

import (
	"context"
	"encoding/json"
	"net/http"
)

//go:generate mockgen -source=main.go -destination=mock_test.go -package=main
type Service interface {
	GetData(ctx context.Context, id string) (any, error)
}

type Handler struct {
	Svc Service
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")

	data, err := h.Svc.GetData(ctx, id)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request cancelled", http.StatusRequestTimeout)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
