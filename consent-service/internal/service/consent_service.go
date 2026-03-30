package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"consent-service/internal/repository"
)

// AuditNotifier defines the interface for sending audit events.
type AuditNotifier interface {
	Notify(ctx context.Context, event AuditEvent) error
}

// AuditEvent represents an event sent to the audit-service.
type AuditEvent struct {
	EventType   string `json:"event_type"`
	EntityType  string `json:"entity_type"`
	EntityID    uint64 `json:"entity_id"`
	PayloadJSON string `json:"payload_json"`
}

// ConsentService defines the business operations for consents.
type ConsentService interface {
	ListConsents(ctx context.Context, limit int) ([]repository.Consent, error)
	CreateConsent(ctx context.Context, c repository.Consent) (uint64, error)
	RevokeConsent(ctx context.Context, documentID string) error
}

type consentService struct {
	repo     repository.ConsentRepository
	notifier AuditNotifier
}

// NewConsentService creates a new ConsentService with the default HTTP audit notifier.
func NewConsentService(repo repository.ConsentRepository) ConsentService {
	auditURL := os.Getenv("AUDIT_SERVICE_URL")
	if auditURL == "" {
		auditURL = "http://localhost:8083"
	}
	return &consentService{
		repo:     repo,
		notifier: &httpAuditNotifier{baseURL: auditURL, client: &http.Client{Timeout: 3 * time.Second}},
	}
}

// NewConsentServiceWithNotifier creates a ConsentService with a custom notifier (for testing).
func NewConsentServiceWithNotifier(repo repository.ConsentRepository, notifier AuditNotifier) ConsentService {
	return &consentService{repo: repo, notifier: notifier}
}

func (s *consentService) ListConsents(ctx context.Context, limit int) ([]repository.Consent, error) {
	return s.repo.List(ctx, limit)
}

func (s *consentService) CreateConsent(ctx context.Context, c repository.Consent) (uint64, error) {
	c.Status = "active"
	id, err := s.repo.Create(ctx, c)
	if err != nil {
		return 0, err
	}

	_ = s.notifier.Notify(ctx, AuditEvent{
		EventType:   "ConsentCreated",
		EntityType:  "consent",
		EntityID:    id,
		PayloadJSON: fmt.Sprintf(`{"consent_id":%d}`, id),
	})

	return id, nil
}

func (s *consentService) RevokeConsent(ctx context.Context, documentID string) error {
	if err := s.repo.Revoke(ctx, documentID); err != nil {
		return err
	}

	_ = s.notifier.Notify(ctx, AuditEvent{
		EventType:   "ConsentRevoked",
		EntityType:  "consent",
		EntityID:    0,
		PayloadJSON: fmt.Sprintf(`{"document_id":"%s"}`, documentID),
	})

	return nil
}

// httpAuditNotifier sends audit events via HTTP POST to the audit-service.
type httpAuditNotifier struct {
	baseURL string
	client  *http.Client
}

func (n *httpAuditNotifier) Notify(ctx context.Context, event AuditEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.baseURL+"/audit-events", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	return nil
}
