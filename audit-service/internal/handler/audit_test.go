package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"audit-service/internal/repository"
)

type mockAuditService struct {
	listFn   func(ctx context.Context, limit int) ([]repository.AuditEvent, error)
	recordFn func(ctx context.Context, e repository.AuditEvent) (uint64, error)
}

func (m *mockAuditService) ListEvents(ctx context.Context, limit int) ([]repository.AuditEvent, error) {
	return m.listFn(ctx, limit)
}

func (m *mockAuditService) RecordEvent(ctx context.Context, e repository.AuditEvent) (uint64, error) {
	return m.recordFn(ctx, e)
}

func TestList_Success(t *testing.T) {
	svc := &mockAuditService{
		listFn: func(_ context.Context, _ int) ([]repository.AuditEvent, error) {
			return []repository.AuditEvent{
				{ID: 1, EventType: "ConsentCreated", EntityType: "consent", EntityID: 10, Payload: "{}"},
			}, nil
		},
	}
	h := NewAuditHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/audit/events", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp []AuditEvent
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp))
	}
}

func TestList_InvalidLimit(t *testing.T) {
	svc := &mockAuditService{}
	h := NewAuditHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/audit/events?limit=abc", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestList_ServiceError(t *testing.T) {
	svc := &mockAuditService{
		listFn: func(_ context.Context, _ int) ([]repository.AuditEvent, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewAuditHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/audit/events", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestList_LimitTooHigh(t *testing.T) {
	svc := &mockAuditService{}
	h := NewAuditHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/audit/events?limit=5000", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecord_Success(t *testing.T) {
	svc := &mockAuditService{
		recordFn: func(_ context.Context, e repository.AuditEvent) (uint64, error) {
			if e.EventType != "ConsentCreated" {
				t.Errorf("expected ConsentCreated, got %s", e.EventType)
			}
			return 99, nil
		},
	}
	h := NewAuditHandler(svc)

	body := `{"event_type":"ConsentCreated","entity_type":"consent","entity_id":10,"payload_json":"{}"}`
	req := httptest.NewRequest(http.MethodPost, "/audit-events", bytes.NewBufferString(body))
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

func TestRecord_InvalidBody(t *testing.T) {
	svc := &mockAuditService{}
	h := NewAuditHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/audit-events", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecord_MissingFields(t *testing.T) {
	svc := &mockAuditService{}
	h := NewAuditHandler(svc)

	body := `{"entity_id":10}`
	req := httptest.NewRequest(http.MethodPost, "/audit-events", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecord_ServiceError(t *testing.T) {
	svc := &mockAuditService{
		recordFn: func(_ context.Context, _ repository.AuditEvent) (uint64, error) {
			return 0, errors.New("db error")
		},
	}
	h := NewAuditHandler(svc)

	body := `{"event_type":"Test","entity_type":"test","entity_id":1,"payload_json":"{}"}`
	req := httptest.NewRequest(http.MethodPost, "/audit-events", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Record(w, req)

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
	svc := &mockAuditService{
		listFn: func(_ context.Context, limit int) ([]repository.AuditEvent, error) {
			if limit != 50 {
				t.Errorf("expected limit=50, got %d", limit)
			}
			return []repository.AuditEvent{}, nil
		},
	}
	h := NewAuditHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/audit/events?limit=50", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
