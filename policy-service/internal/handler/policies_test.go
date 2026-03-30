package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"policy-service/internal/repository"
)

type mockPolicyService struct {
	listFn   func(ctx context.Context, limit int) ([]repository.Policy, error)
	createFn func(ctx context.Context, p repository.Policy) (uint64, error)
}

func (m *mockPolicyService) ListPolicies(ctx context.Context, limit int) ([]repository.Policy, error) {
	return m.listFn(ctx, limit)
}

func (m *mockPolicyService) CreatePolicy(ctx context.Context, p repository.Policy) (uint64, error) {
	return m.createFn(ctx, p)
}

func TestList_Success(t *testing.T) {
	svc := &mockPolicyService{
		listFn: func(_ context.Context, _ int) ([]repository.Policy, error) {
			return []repository.Policy{
				{ID: 1, Version: "1.0", ContentHash: "abc123"},
			}, nil
		},
	}
	h := NewPolicyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/policies", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []Policy
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(resp))
	}
}

func TestList_InvalidLimit(t *testing.T) {
	svc := &mockPolicyService{}
	h := NewPolicyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/policies?limit=abc", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestList_ServiceError(t *testing.T) {
	svc := &mockPolicyService{
		listFn: func(_ context.Context, _ int) ([]repository.Policy, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewPolicyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/policies", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestList_LimitTooHigh(t *testing.T) {
	svc := &mockPolicyService{}
	h := NewPolicyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/policies?limit=5000", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreate_Success(t *testing.T) {
	svc := &mockPolicyService{
		createFn: func(_ context.Context, p repository.Policy) (uint64, error) {
			if p.Version != "2.0" {
				t.Errorf("expected version=2.0, got %s", p.Version)
			}
			return 42, nil
		},
	}
	h := NewPolicyHandler(svc)

	body := `{"version":"2.0","content_hash":"xyz789"}`
	req := httptest.NewRequest(http.MethodPost, "/policies", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var resp map[string]uint64
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["id"] != 42 {
		t.Errorf("expected id=42, got %d", resp["id"])
	}
}

func TestCreate_InvalidBody(t *testing.T) {
	svc := &mockPolicyService{}
	h := NewPolicyHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/policies", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreate_MissingFields(t *testing.T) {
	svc := &mockPolicyService{}
	h := NewPolicyHandler(svc)

	body := `{"version":"2.0"}`
	req := httptest.NewRequest(http.MethodPost, "/policies", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreate_ServiceError(t *testing.T) {
	svc := &mockPolicyService{
		createFn: func(_ context.Context, _ repository.Policy) (uint64, error) {
			return 0, errors.New("db error")
		},
	}
	h := NewPolicyHandler(svc)

	body := `{"version":"2.0","content_hash":"xyz789"}`
	req := httptest.NewRequest(http.MethodPost, "/policies", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
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

func TestList_CustomLimit(t *testing.T) {
	svc := &mockPolicyService{
		listFn: func(_ context.Context, limit int) ([]repository.Policy, error) {
			if limit != 50 {
				t.Errorf("expected limit=50, got %d", limit)
			}
			return []repository.Policy{}, nil
		},
	}
	h := NewPolicyHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/policies?limit=50", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
