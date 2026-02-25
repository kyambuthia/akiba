package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "PORT") {
		t.Fatalf("expected PORT parse error, got %v", err)
	}
}

func TestLoadRejectsInvalidDuration(t *testing.T) {
	t.Setenv("ACCESS_TOKEN_TTL", "abc")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "ACCESS_TOKEN_TTL") {
		t.Fatalf("expected ACCESS_TOKEN_TTL parse error, got %v", err)
	}
}

func TestLoadRejectsNonPositiveDurations(t *testing.T) {
	t.Setenv("DB_TIMEOUT", "0s")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "DB_TIMEOUT") {
		t.Fatalf("expected DB_TIMEOUT validation error, got %v", err)
	}
}

func TestLoadDefaultsRemainValid(t *testing.T) {
	for _, k := range []string{"PORT", "ACCESS_TOKEN_TTL", "DB_TIMEOUT"} {
		_ = os.Unsetenv(k)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected defaults to load: %v", err)
	}
	if cfg.Port != 8080 {
		t.Fatalf("unexpected default port: %d", cfg.Port)
	}
}
