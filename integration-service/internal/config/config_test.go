package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	for _, key := range []string{"APP_ADDR", "EXTERNAL_BASE_URL"} {
		os.Unsetenv(key)
	}

	cfg := Load()

	if cfg.Addr != ":8085" {
		t.Errorf("expected :8085, got %s", cfg.Addr)
	}
	if cfg.ExternalBase != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", cfg.ExternalBase)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("APP_ADDR", ":9090")
	os.Setenv("EXTERNAL_BASE_URL", "https://custom.com")
	defer func() {
		os.Unsetenv("APP_ADDR")
		os.Unsetenv("EXTERNAL_BASE_URL")
	}()

	cfg := Load()

	if cfg.Addr != ":9090" {
		t.Errorf("expected :9090, got %s", cfg.Addr)
	}
	if cfg.ExternalBase != "https://custom.com" {
		t.Errorf("expected https://custom.com, got %s", cfg.ExternalBase)
	}
}

func TestGetEnv_Fallback(t *testing.T) {
	os.Unsetenv("NONEXISTENT_KEY")
	v := getEnv("NONEXISTENT_KEY", "default")
	if v != "default" {
		t.Errorf("expected default, got %s", v)
	}
}

func TestGetEnv_Set(t *testing.T) {
	os.Setenv("TEST_KEY_INTEGRATION", "value123")
	defer os.Unsetenv("TEST_KEY_INTEGRATION")

	v := getEnv("TEST_KEY_INTEGRATION", "default")
	if v != "value123" {
		t.Errorf("expected value123, got %s", v)
	}
}
