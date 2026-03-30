package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"consent-service/internal/repository"
)

type mockRepo struct {
	listFn   func(ctx context.Context, limit int) ([]repository.Consent, error)
	createFn func(ctx context.Context, c repository.Consent) (uint64, error)
	revokeFn func(ctx context.Context, documentID string) error
}

func (m *mockRepo) List(ctx context.Context, limit int) ([]repository.Consent, error) {
	return m.listFn(ctx, limit)
}

func (m *mockRepo) Create(ctx context.Context, c repository.Consent) (uint64, error) {
	return m.createFn(ctx, c)
}

func (m *mockRepo) Revoke(ctx context.Context, documentID string) error {
	return m.revokeFn(ctx, documentID)
}

type mockNotifier struct {
	notifyFn func(ctx context.Context, event AuditEvent) error
}

func (m *mockNotifier) Notify(ctx context.Context, event AuditEvent) error {
	if m.notifyFn != nil {
		return m.notifyFn(ctx, event)
	}
	return nil
}

func TestListConsents_Success(t *testing.T) {
	repo := &mockRepo{
		listFn: func(_ context.Context, limit int) ([]repository.Consent, error) {
			return []repository.Consent{{ID: 1}}, nil
		},
	}
	svc := NewConsentServiceWithNotifier(repo, &mockNotifier{})

	result, err := svc.ListConsents(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 consent, got %d", len(result))
	}
}

func TestListConsents_Error(t *testing.T) {
	repo := &mockRepo{
		listFn: func(_ context.Context, _ int) ([]repository.Consent, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewConsentServiceWithNotifier(repo, &mockNotifier{})

	_, err := svc.ListConsents(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateConsent_Success(t *testing.T) {
	var notified bool
	repo := &mockRepo{
		createFn: func(_ context.Context, c repository.Consent) (uint64, error) {
			if c.Status != "active" {
				t.Errorf("expected status=active, got %s", c.Status)
			}
			return 42, nil
		},
	}
	notifier := &mockNotifier{
		notifyFn: func(_ context.Context, event AuditEvent) error {
			notified = true
			if event.EventType != "ConsentCreated" {
				t.Errorf("expected ConsentCreated, got %s", event.EventType)
			}
			if event.EntityID != 42 {
				t.Errorf("expected entity_id=42, got %d", event.EntityID)
			}
			return nil
		},
	}
	svc := NewConsentServiceWithNotifier(repo, notifier)

	id, err := svc.CreateConsent(context.Background(), repository.Consent{UserID: 1, PolicyID: 2, Purpose: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id=42, got %d", id)
	}
	if !notified {
		t.Error("expected notifier to be called")
	}
}

func TestCreateConsent_RepoError(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, _ repository.Consent) (uint64, error) {
			return 0, errors.New("insert error")
		},
	}
	svc := NewConsentServiceWithNotifier(repo, &mockNotifier{})

	_, err := svc.CreateConsent(context.Background(), repository.Consent{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRevokeConsent_Success(t *testing.T) {
	var notified bool
	repo := &mockRepo{
		revokeFn: func(_ context.Context, docID string) error {
			if docID != "doc-abc" {
				t.Errorf("expected doc-abc, got %s", docID)
			}
			return nil
		},
	}
	notifier := &mockNotifier{
		notifyFn: func(_ context.Context, event AuditEvent) error {
			notified = true
			if event.EventType != "ConsentRevoked" {
				t.Errorf("expected ConsentRevoked, got %s", event.EventType)
			}
			return nil
		},
	}
	svc := NewConsentServiceWithNotifier(repo, notifier)

	err := svc.RevokeConsent(context.Background(), "doc-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !notified {
		t.Error("expected notifier to be called")
	}
}

func TestRevokeConsent_RepoError(t *testing.T) {
	repo := &mockRepo{
		revokeFn: func(_ context.Context, _ string) error {
			return errors.New("update error")
		},
	}
	svc := NewConsentServiceWithNotifier(repo, &mockNotifier{})

	err := svc.RevokeConsent(context.Background(), "doc-abc")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewConsentService_DefaultNotifier(t *testing.T) {
	repo := &mockRepo{}
	svc := NewConsentService(repo)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestNewConsentService_WithEnvVar(t *testing.T) {
	t.Setenv("AUDIT_SERVICE_URL", "http://custom:9999")
	repo := &mockRepo{}
	svc := NewConsentService(repo)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestHTTPAuditNotifier_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/audit-events" {
			t.Errorf("expected /audit-events, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	notifier := &httpAuditNotifier{
		baseURL: server.URL,
		client:  server.Client(),
	}

	err := notifier.Notify(context.Background(), AuditEvent{
		EventType:   "Test",
		EntityType:  "test",
		EntityID:    1,
		PayloadJSON: "{}",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTPAuditNotifier_ServerError(t *testing.T) {
	notifier := &httpAuditNotifier{
		baseURL: "http://localhost:1",
		client:  &http.Client{Timeout: 1 * time.Second},
	}

	err := notifier.Notify(context.Background(), AuditEvent{
		EventType:   "Test",
		EntityType:  "test",
		EntityID:    1,
		PayloadJSON: "{}",
	})
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
