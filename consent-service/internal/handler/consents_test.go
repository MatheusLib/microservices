package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"consent-service/internal/repository"
)

// mockConsentService implements service.ConsentService for testing.
type mockConsentService struct {
	listFn   func(ctx context.Context, limit int) ([]repository.Consent, error)
	createFn func(ctx context.Context, c repository.Consent) (uint64, error)
	revokeFn func(ctx context.Context, documentID string) error
}

func (m *mockConsentService) ListConsents(ctx context.Context, limit int) ([]repository.Consent, error) {
	return m.listFn(ctx, limit)
}

func (m *mockConsentService) CreateConsent(ctx context.Context, c repository.Consent) (uint64, error) {
	return m.createFn(ctx, c)
}

func (m *mockConsentService) RevokeConsent(ctx context.Context, documentID string) error {
	return m.revokeFn(ctx, documentID)
}

func TestList_Success(t *testing.T) {
	svc := &mockConsentService{
		listFn: func(_ context.Context, _ int) ([]repository.Consent, error) {
			return []repository.Consent{
				{ID: 1, UserID: 10, PolicyID: 20, Purpose: "marketing", Status: "active"},
			}, nil
		},
	}
	h := NewConsentHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/consents", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []Consent
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 consent, got %d", len(resp))
	}
	if resp[0].ID != 1 {
		t.Errorf("expected id=1, got %d", resp[0].ID)
	}
}

func TestList_InvalidLimit(t *testing.T) {
	svc := &mockConsentService{}
	h := NewConsentHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/consents?limit=abc", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestList_ServiceError(t *testing.T) {
	svc := &mockConsentService{
		listFn: func(_ context.Context, _ int) ([]repository.Consent, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewConsentHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/consents", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCreate_Success(t *testing.T) {
	svc := &mockConsentService{
		createFn: func(_ context.Context, c repository.Consent) (uint64, error) {
			return 42, nil
		},
	}
	h := NewConsentHandler(svc)

	body := `{"user_id":10,"policy_id":20,"purpose":"analytics"}`
	req := httptest.NewRequest(http.MethodPost, "/consents", bytes.NewBufferString(body))
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
	svc := &mockConsentService{}
	h := NewConsentHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/consents", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreate_MissingFields(t *testing.T) {
	svc := &mockConsentService{}
	h := NewConsentHandler(svc)

	body := `{"user_id":10}`
	req := httptest.NewRequest(http.MethodPost, "/consents", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreate_ServiceError(t *testing.T) {
	svc := &mockConsentService{
		createFn: func(_ context.Context, _ repository.Consent) (uint64, error) {
			return 0, errors.New("db error")
		},
	}
	h := NewConsentHandler(svc)

	body := `{"user_id":10,"policy_id":20,"purpose":"analytics"}`
	req := httptest.NewRequest(http.MethodPost, "/consents", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestRevoke_Success(t *testing.T) {
	svc := &mockConsentService{
		revokeFn: func(_ context.Context, docID string) error {
			if docID != "doc123" {
				t.Errorf("expected doc123, got %s", docID)
			}
			return nil
		},
	}
	h := NewConsentHandler(svc)

	r := chi.NewRouter()
	r.Patch("/consents/{document_id}/revoke", h.Revoke)

	req := httptest.NewRequest(http.MethodPatch, "/consents/doc123/revoke", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestRevoke_ServiceError(t *testing.T) {
	svc := &mockConsentService{
		revokeFn: func(_ context.Context, _ string) error {
			return errors.New("db error")
		},
	}
	h := NewConsentHandler(svc)

	r := chi.NewRouter()
	r.Patch("/consents/{document_id}/revoke", h.Revoke)

	req := httptest.NewRequest(http.MethodPatch, "/consents/doc123/revoke", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

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
	svc := &mockConsentService{
		listFn: func(_ context.Context, limit int) ([]repository.Consent, error) {
			if limit != 50 {
				t.Errorf("expected limit=50, got %d", limit)
			}
			return []repository.Consent{}, nil
		},
	}
	h := NewConsentHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/consents?limit=50", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestList_LimitTooHigh(t *testing.T) {
	svc := &mockConsentService{}
	h := NewConsentHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/consents?limit=5000", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
