package service

import (
	"context"
	"errors"
	"testing"

	"audit-service/internal/repository"
)

type mockRepo struct {
	listFn   func(ctx context.Context, limit int) ([]repository.AuditEvent, error)
	recordFn func(ctx context.Context, e repository.AuditEvent) (uint64, error)
}

func (m *mockRepo) List(ctx context.Context, limit int) ([]repository.AuditEvent, error) {
	return m.listFn(ctx, limit)
}

func (m *mockRepo) Record(ctx context.Context, e repository.AuditEvent) (uint64, error) {
	return m.recordFn(ctx, e)
}

func TestListEvents_Success(t *testing.T) {
	repo := &mockRepo{
		listFn: func(_ context.Context, _ int) ([]repository.AuditEvent, error) {
			return []repository.AuditEvent{{ID: 1}}, nil
		},
	}
	svc := NewAuditService(repo)

	result, err := svc.ListEvents(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 event, got %d", len(result))
	}
}

func TestListEvents_Error(t *testing.T) {
	repo := &mockRepo{
		listFn: func(_ context.Context, _ int) ([]repository.AuditEvent, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewAuditService(repo)

	_, err := svc.ListEvents(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecordEvent_Success(t *testing.T) {
	repo := &mockRepo{
		recordFn: func(_ context.Context, e repository.AuditEvent) (uint64, error) {
			if e.EventType != "Test" {
				t.Errorf("expected Test, got %s", e.EventType)
			}
			return 42, nil
		},
	}
	svc := NewAuditService(repo)

	id, err := svc.RecordEvent(context.Background(), repository.AuditEvent{EventType: "Test", EntityType: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id=42, got %d", id)
	}
}

func TestRecordEvent_Error(t *testing.T) {
	repo := &mockRepo{
		recordFn: func(_ context.Context, _ repository.AuditEvent) (uint64, error) {
			return 0, errors.New("insert error")
		},
	}
	svc := NewAuditService(repo)

	_, err := svc.RecordEvent(context.Background(), repository.AuditEvent{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
