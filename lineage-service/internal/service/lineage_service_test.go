package service

import (
	"context"
	"errors"
	"testing"

	"lineage-service/internal/repository"
)

type mockRepo struct {
	recordFn       func(ctx context.Context, e repository.LineageEvent) (uint64, error)
	listBySubjFn   func(ctx context.Context, subjectID uint64) ([]repository.LineageEvent, error)
}

func (m *mockRepo) Record(ctx context.Context, e repository.LineageEvent) (uint64, error) {
	return m.recordFn(ctx, e)
}

func (m *mockRepo) ListBySubject(ctx context.Context, subjectID uint64) ([]repository.LineageEvent, error) {
	return m.listBySubjFn(ctx, subjectID)
}

func TestRecord_Success(t *testing.T) {
	repo := &mockRepo{
		recordFn: func(_ context.Context, e repository.LineageEvent) (uint64, error) {
			if e.Operation != "COLLECT" {
				t.Errorf("expected COLLECT, got %s", e.Operation)
			}
			return 42, nil
		},
	}
	svc := NewLineageService(repo)

	id, err := svc.Record(context.Background(), repository.LineageEvent{
		SubjectID: 1, Operation: "COLLECT", Source: "api", Destination: "db", Purpose: "analytics",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id=42, got %d", id)
	}
}

func TestRecord_Error(t *testing.T) {
	repo := &mockRepo{
		recordFn: func(_ context.Context, _ repository.LineageEvent) (uint64, error) {
			return 0, errors.New("insert error")
		},
	}
	svc := NewLineageService(repo)

	_, err := svc.Record(context.Background(), repository.LineageEvent{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExportBySubject_Success(t *testing.T) {
	repo := &mockRepo{
		listBySubjFn: func(_ context.Context, subjectID uint64) ([]repository.LineageEvent, error) {
			if subjectID != 100 {
				t.Errorf("expected subjectID=100, got %d", subjectID)
			}
			return []repository.LineageEvent{
				{ID: 1, SubjectID: 100, Operation: "COLLECT"},
				{ID: 2, SubjectID: 100, Operation: "SHARE"},
			}, nil
		},
	}
	svc := NewLineageService(repo)

	result, err := svc.ExportBySubject(context.Background(), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 events, got %d", len(result))
	}
}

func TestExportBySubject_Error(t *testing.T) {
	repo := &mockRepo{
		listBySubjFn: func(_ context.Context, _ uint64) ([]repository.LineageEvent, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewLineageService(repo)

	_, err := svc.ExportBySubject(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
