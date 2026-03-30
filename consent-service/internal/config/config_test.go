package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Ensure env vars are not set
	for _, key := range []string{"APP_ADDR", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME"} {
		os.Unsetenv(key)
	}

	cfg := Load()

	if cfg.Addr != ":8081" {
		t.Errorf("expected :8081, got %s", cfg.Addr)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "3306" {
		t.Errorf("expected 3306, got %s", cfg.DBPort)
	}
	if cfg.DBUser != "admin" {
		t.Errorf("expected admin, got %s", cfg.DBUser)
	}
	if cfg.DBName != "tcc" {
		t.Errorf("expected tcc, got %s", cfg.DBName)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("APP_ADDR", ":9090")
	os.Setenv("DB_HOST", "myhost")
	defer func() {
		os.Unsetenv("APP_ADDR")
		os.Unsetenv("DB_HOST")
	}()

	cfg := Load()

	if cfg.Addr != ":9090" {
		t.Errorf("expected :9090, got %s", cfg.Addr)
	}
	if cfg.DBHost != "myhost" {
		t.Errorf("expected myhost, got %s", cfg.DBHost)
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
	os.Setenv("TEST_KEY_CONSENT", "value123")
	defer os.Unsetenv("TEST_KEY_CONSENT")

	v := getEnv("TEST_KEY_CONSENT", "default")
	if v != "value123" {
		t.Errorf("expected value123, got %s", v)
	}
}
