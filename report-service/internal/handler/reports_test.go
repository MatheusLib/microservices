package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"report-service/internal/repository"
)

type mockReportService struct {
	listConsentsFn func(ctx context.Context, userID *uint64, limit int) ([]repository.ConsentReport, error)
}

func (m *mockReportService) ListConsents(ctx context.Context, userID *uint64, limit int) ([]repository.ConsentReport, error) {
	return m.listConsentsFn(ctx, userID, limit)
}

func TestListConsents_Success(t *testing.T) {
	svc := &mockReportService{
		listConsentsFn: func(_ context.Context, _ *uint64, _ int) ([]repository.ConsentReport, error) {
			return []repository.ConsentReport{
				{ID: 1, UserID: 10, PolicyID: 20, Purpose: "marketing", Status: "active"},
			}, nil
		},
	}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []ConsentReport
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 report, got %d", len(resp))
	}
}

func TestListConsents_WithUserID(t *testing.T) {
	svc := &mockReportService{
		listConsentsFn: func(_ context.Context, userID *uint64, _ int) ([]repository.ConsentReport, error) {
			if userID == nil || *userID != 42 {
				t.Error("expected userID=42")
			}
			return []repository.ConsentReport{}, nil
		},
	}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents?user_id=42", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestListConsents_InvalidUserID(t *testing.T) {
	svc := &mockReportService{}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents?user_id=abc", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListConsents_InvalidLimit(t *testing.T) {
	svc := &mockReportService{}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents?limit=abc", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListConsents_LimitTooHigh(t *testing.T) {
	svc := &mockReportService{}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents?limit=5000", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListConsents_ServiceError(t *testing.T) {
	svc := &mockReportService{
		listConsentsFn: func(_ context.Context, _ *uint64, _ int) ([]repository.ConsentReport, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

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

func TestListConsents_CustomLimit(t *testing.T) {
	svc := &mockReportService{
		listConsentsFn: func(_ context.Context, _ *uint64, limit int) ([]repository.ConsentReport, error) {
			if limit != 50 {
				t.Errorf("expected limit=50, got %d", limit)
			}
			return []repository.ConsentReport{}, nil
		},
	}
	h := NewReportHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/reports/consents?limit=50", nil)
	w := httptest.NewRecorder()
	h.ListConsents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
