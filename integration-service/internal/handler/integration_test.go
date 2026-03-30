package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
	pingFn func(ctx context.Context) error
}

func (m *mockService) Ping(ctx context.Context) error {
	return m.pingFn(ctx)
}

func TestPing_Success(t *testing.T) {
	svc := &mockService{
		pingFn: func(_ context.Context) error { return nil },
	}
	h := NewIntegrationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/integrations/ping", nil)
	w := httptest.NewRecorder()
	h.Ping(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected 'ok', got %s", w.Body.String())
	}
}

func TestPing_Error(t *testing.T) {
	svc := &mockService{
		pingFn: func(_ context.Context) error { return errors.New("unreachable") },
	}
	h := NewIntegrationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/integrations/ping", nil)
	w := httptest.NewRecorder()
	h.Ping(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected 'ok', got %s", w.Body.String())
	}
}
