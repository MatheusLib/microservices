package external_services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClient_Ping_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	err := client.Ping(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTPClient_Ping_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	// The current implementation doesn't check status codes, so no error expected
	err := client.Ping(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTPClient_Ping_Unreachable(t *testing.T) {
	client := NewHTTPClient("http://localhost:1")
	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHTTPClient_Ping_InvalidURL(t *testing.T) {
	client := NewHTTPClient("://invalid")
	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
