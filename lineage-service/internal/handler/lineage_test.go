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

	"lineage-service/internal/repository"
)

type mockLineageService struct {
	recordFn       func(ctx context.Context, e repository.LineageEvent) (uint64, error)
	exportBySubjFn func(ctx context.Context, subjectID uint64) ([]repository.LineageEvent, error)
}

func (m *mockLineageService) Record(ctx context.Context, e repository.LineageEvent) (uint64, error) {
	return m.recordFn(ctx, e)
}

func (m *mockLineageService) ExportBySubject(ctx context.Context, subjectID uint64) ([]repository.LineageEvent, error) {
	return m.exportBySubjFn(ctx, subjectID)
}

func TestRecord_Success(t *testing.T) {
	svc := &mockLineageService{
		recordFn: func(_ context.Context, e repository.LineageEvent) (uint64, error) {
			if e.Operation != "COLLECT" {
				t.Errorf("expected COLLECT, got %s", e.Operation)
			}
			return 99, nil
		},
	}
	h := NewLineageHandler(svc)

	body := `{"subject_id":1,"operation":"COLLECT","source":"api","destination":"db","purpose":"analytics"}`
	req := httptest.NewRequest(http.MethodPost, "/lineage", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var resp map[string]uint64
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["id"] != 99 {
		t.Errorf("expected id=99, got %d", resp["id"])
	}
}

func TestRecord_WithConsentID(t *testing.T) {
	svc := &mockLineageService{
		recordFn: func(_ context.Context, e repository.LineageEvent) (uint64, error) {
			if e.ConsentID == nil || *e.ConsentID != 42 {
				t.Errorf("expected consent_id=42")
			}
			return 10, nil
		},
	}
	h := NewLineageHandler(svc)

	body := `{"subject_id":1,"operation":"COLLECT","source":"api","destination":"db","purpose":"analytics","consent_id":42}`
	req := httptest.NewRequest(http.MethodPost, "/lineage", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestRecord_InvalidBody(t *testing.T) {
	svc := &mockLineageService{}
	h := NewLineageHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/lineage", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecord_MissingFields(t *testing.T) {
	svc := &mockLineageService{}
	h := NewLineageHandler(svc)

	body := `{"subject_id":1,"operation":"COLLECT"}`
	req := httptest.NewRequest(http.MethodPost, "/lineage", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecord_MissingSubjectID(t *testing.T) {
	svc := &mockLineageService{}
	h := NewLineageHandler(svc)

	body := `{"operation":"COLLECT","source":"api","destination":"db","purpose":"analytics"}`
	req := httptest.NewRequest(http.MethodPost, "/lineage", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecord_ServiceError(t *testing.T) {
	svc := &mockLineageService{
		recordFn: func(_ context.Context, _ repository.LineageEvent) (uint64, error) {
			return 0, errors.New("db error")
		},
	}
	h := NewLineageHandler(svc)

	body := `{"subject_id":1,"operation":"COLLECT","source":"api","destination":"db","purpose":"analytics"}`
	req := httptest.NewRequest(http.MethodPost, "/lineage", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestExport_Success(t *testing.T) {
	consentID := uint64(10)
	svc := &mockLineageService{
		exportBySubjFn: func(_ context.Context, subjectID uint64) ([]repository.LineageEvent, error) {
			if subjectID != 100 {
				t.Errorf("expected subjectID=100, got %d", subjectID)
			}
			return []repository.LineageEvent{
				{ID: 1, SubjectID: 100, Operation: "COLLECT", Source: "api", Destination: "db", Purpose: "analytics", ConsentID: &consentID, PayloadJSON: "{}", CreatedAt: "2025-01-01"},
			}, nil
		},
	}
	h := NewLineageHandler(svc)

	r := chi.NewRouter()
	r.Get("/lineage/export/{subject_id}", h.Export)

	req := httptest.NewRequest(http.MethodGet, "/lineage/export/100", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []LineageEventResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp))
	}
	if resp[0].Operation != "COLLECT" {
		t.Errorf("expected COLLECT, got %s", resp[0].Operation)
	}
}

func TestExport_InvalidSubjectID(t *testing.T) {
	svc := &mockLineageService{}
	h := NewLineageHandler(svc)

	r := chi.NewRouter()
	r.Get("/lineage/export/{subject_id}", h.Export)

	req := httptest.NewRequest(http.MethodGet, "/lineage/export/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExport_ZeroSubjectID(t *testing.T) {
	svc := &mockLineageService{}
	h := NewLineageHandler(svc)

	r := chi.NewRouter()
	r.Get("/lineage/export/{subject_id}", h.Export)

	req := httptest.NewRequest(http.MethodGet, "/lineage/export/0", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExport_ServiceError(t *testing.T) {
	svc := &mockLineageService{
		exportBySubjFn: func(_ context.Context, _ uint64) ([]repository.LineageEvent, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewLineageHandler(svc)

	r := chi.NewRouter()
	r.Get("/lineage/export/{subject_id}", h.Export)

	req := httptest.NewRequest(http.MethodGet, "/lineage/export/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestExport_Empty(t *testing.T) {
	svc := &mockLineageService{
		exportBySubjFn: func(_ context.Context, _ uint64) ([]repository.LineageEvent, error) {
			return []repository.LineageEvent{}, nil
		},
	}
	h := NewLineageHandler(svc)

	r := chi.NewRouter()
	r.Get("/lineage/export/{subject_id}", h.Export)

	req := httptest.NewRequest(http.MethodGet, "/lineage/export/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []LineageEventResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 0 {
		t.Fatalf("expected 0 events, got %d", len(resp))
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
