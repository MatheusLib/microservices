package external_services

import (
	"context"
	"errors"
	"testing"
)

type mockClient struct {
	pingFn func(ctx context.Context) error
}

func (m *mockClient) Ping(ctx context.Context) error {
	return m.pingFn(ctx)
}

func TestService_Ping_Success(t *testing.T) {
	client := &mockClient{
		pingFn: func(_ context.Context) error { return nil },
	}
	svc := NewService(client)

	err := svc.Ping(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_Ping_Error(t *testing.T) {
	client := &mockClient{
		pingFn: func(_ context.Context) error { return errors.New("fail") },
	}
	svc := NewService(client)

	err := svc.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
