package cli

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected default format json, got %q", cfg.Format)
	}
	if cfg.ShowStats {
		t.Error("expected ShowStats false by default")
	}
}

func TestParseFlags_FullArgs(t *testing.T) {
	cfg, err := parseFlags([]string{
		"-input", "app.log",
		"-output", "out.log",
		"-format", "text",
		"-from", "2024-01-01T00:00:00Z",
		"-to", "2024-01-02T00:00:00Z",
		"-field", "level=error",
		"-field", "service=api",
		"-stats",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Input != "app.log" {
		t.Errorf("expected input app.log, got %q", cfg.Input)
	}
	if cfg.Format != "text" {
		t.Errorf("expected format text, got %q", cfg.Format)
	}
	if cfg.Fields["level"] != "error" {
		t.Errorf("expected level=error, got %q", cfg.Fields["level"])
	}
	if cfg.Fields["service"] != "api" {
		t.Errorf("expected service=api, got %q", cfg.Fields["service"])
	}
	if !cfg.ShowStats {
		t.Error("expected ShowStats true")
	}
}

func TestParseFlags_InvalidFieldFormat(t *testing.T) {
	_, err := parseFlags([]string{"-field", "badvalue"})
	if err == nil {
		t.Error("expected error for invalid field format")
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := parseFlags([]string{"-format", "xml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
