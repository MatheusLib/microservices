package service

import (
	"context"
	"errors"
	"testing"

	"report-service/internal/repository"
)

type mockRepo struct {
	listConsentsFn func(ctx context.Context, userID *uint64, limit int) ([]repository.ConsentReport, error)
}

func (m *mockRepo) ListConsents(ctx context.Context, userID *uint64, limit int) ([]repository.ConsentReport, error) {
	return m.listConsentsFn(ctx, userID, limit)
}

func TestListConsents_Success(t *testing.T) {
	repo := &mockRepo{
		listConsentsFn: func(_ context.Context, _ *uint64, _ int) ([]repository.ConsentReport, error) {
			return []repository.ConsentReport{{ID: 1}}, nil
		},
	}
	svc := NewReportService(repo)

	result, err := svc.ListConsents(context.Background(), nil, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 report, got %d", len(result))
	}
}

func TestListConsents_WithUserID(t *testing.T) {
	uid := uint64(42)
	repo := &mockRepo{
		listConsentsFn: func(_ context.Context, userID *uint64, _ int) ([]repository.ConsentReport, error) {
			if userID == nil || *userID != 42 {
				t.Error("expected userID=42")
			}
			return []repository.ConsentReport{}, nil
		},
	}
	svc := NewReportService(repo)

	_, err := svc.ListConsents(context.Background(), &uid, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListConsents_Error(t *testing.T) {
	repo := &mockRepo{
		listConsentsFn: func(_ context.Context, _ *uint64, _ int) ([]repository.ConsentReport, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewReportService(repo)

	_, err := svc.ListConsents(context.Background(), nil, 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
