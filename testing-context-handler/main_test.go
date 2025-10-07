package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestHandler_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	m := NewMockService(ctrl)
	m.EXPECT().GetData(gomock.Any(), "key1").DoAndReturn(
		func(ctx context.Context, id string) (any, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return map[string]any{"id": id, "ok": true}, nil
			}
		},
	)

	h := Handler{Svc: m}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/?id=key1", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("expected status ok")
	}
	expected := `{"id":"key1","ok":true}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestHandler_Cancelled(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	m := NewMockService(ctrl)
	m.EXPECT().GetData(gomock.Any(), "key1").DoAndReturn(
		func(ctx context.Context, id string) (any, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	)

	h := Handler{Svc: m}

	ctx, _ := context.WithTimeout(t.Context(), 50*time.Millisecond)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/?id=key1", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusRequestTimeout {
		t.Error("expected status request timeout")
	}
	expected := "request cancelled"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestHandler_InternalServerError(t *testing.T) {
		t.Parallel()

	ctrl := gomock.NewController(t)

	m := NewMockService(ctrl)
	m.EXPECT().GetData(gomock.Any(), "key1").DoAndReturn(
		func(ctx context.Context, id string) (any, error) {
			return nil, errors.New("failed to get data")
		},
	)

	h := Handler{Svc: m}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/?id=key1", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Error("expected status internal server error")
	}
	expected := "internal error"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body %q, got %q", expected, rr.Body.String())
	}
}
