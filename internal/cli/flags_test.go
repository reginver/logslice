package cli

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "auto" {
		t.Errorf("expected format=auto, got %q", cfg.Format)
	}
	if cfg.OutputFmt != "json" {
		t.Errorf("expected output-format=json, got %q", cfg.OutputFmt)
	}
	if cfg.Stats {
		t.Error("expected stats=false")
	}
}

func TestParseFlags_FullArgs(t *testing.T) {
	cfg, err := parseFlags([]string{
		"-input", "in.log",
		"-output", "out.log",
		"-start", "2024-01-01T00:00:00Z",
		"-end", "2024-01-02T00:00:00Z",
		"-field", "level=error",
		"-regex", "msg=timeout.*",
		"-stats",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Input != "in.log" {
		t.Errorf("input: got %q", cfg.Input)
	}
	if cfg.Fields["level"] != "error" {
		t.Errorf("field level: got %q", cfg.Fields["level"])
	}
	if cfg.Regexes["msg"] != "timeout.*" {
		t.Errorf("regex msg: got %q", cfg.Regexes["msg"])
	}
	if !cfg.Stats {
		t.Error("expected stats=true")
	}
}

func TestParseFlags_InvalidFieldFormat(t *testing.T) {
	_, err := parseFlags([]string{"-field", "badvalue"})
	if err == nil {
		t.Fatal("expected error for invalid field format")
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := parseFlags([]string{"-regex", "=pattern"})
	if err == nil {
		t.Fatal("expected error for empty field name in regex")
	}
}
