package service

import (
	"context"
	"errors"
	"testing"

	"policy-service/internal/repository"
)

type mockRepo struct {
	listFn   func(ctx context.Context, limit int) ([]repository.Policy, error)
	createFn func(ctx context.Context, p repository.Policy) (uint64, error)
}

func (m *mockRepo) List(ctx context.Context, limit int) ([]repository.Policy, error) {
	return m.listFn(ctx, limit)
}

func (m *mockRepo) Create(ctx context.Context, p repository.Policy) (uint64, error) {
	return m.createFn(ctx, p)
}

func TestListPolicies_Success(t *testing.T) {
	repo := &mockRepo{
		listFn: func(_ context.Context, _ int) ([]repository.Policy, error) {
			return []repository.Policy{{ID: 1, Version: "1.0", ContentHash: "abc"}}, nil
		},
	}
	svc := NewPolicyService(repo)

	result, err := svc.ListPolicies(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(result))
	}
}

func TestListPolicies_Error(t *testing.T) {
	repo := &mockRepo{
		listFn: func(_ context.Context, _ int) ([]repository.Policy, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewPolicyService(repo)

	_, err := svc.ListPolicies(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreatePolicy_Success(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, p repository.Policy) (uint64, error) {
			if p.Version != "2.0" {
				t.Errorf("expected version=2.0, got %s", p.Version)
			}
			return 42, nil
		},
	}
	svc := NewPolicyService(repo)

	id, err := svc.CreatePolicy(context.Background(), repository.Policy{Version: "2.0", ContentHash: "xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id=42, got %d", id)
	}
}

func TestCreatePolicy_Error(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, _ repository.Policy) (uint64, error) {
			return 0, errors.New("insert error")
		},
	}
	svc := NewPolicyService(repo)

	_, err := svc.CreatePolicy(context.Background(), repository.Policy{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
