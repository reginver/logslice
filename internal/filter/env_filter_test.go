package filter

import (
	"testing"
)

func envEntry(field, value string) map[string]interface{} {
	return map[string]interface{}{field: value}
}

func TestNewEnvFilter_Valid(t *testing.T) {
	f, err := NewEnvFilter("env", []string{"production", "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewEnvFilter_EmptyField(t *testing.T) {
	_, err := NewEnvFilter("", []string{"production"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEnvFilter_NoEnvs(t *testing.T) {
	_, err := NewEnvFilter("env", []string{})
	if err == nil {
		t.Fatal("expected error for empty env list")
	}
}

func TestNewEnvFilter_EmptyEnvName(t *testing.T) {
	_, err := NewEnvFilter("env", []string{"production", ""})
	if err == nil {
		t.Fatal("expected error for empty env name")
	}
}

func TestEnvFilter_Match_Hit(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"production", "staging"})
	if !f.Match(envEntry("env", "production")) {
		t.Error("expected match for 'production'")
	}
}

func TestEnvFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"Production"})
	if !f.Match(envEntry("env", "PRODUCTION")) {
		t.Error("expected case-insensitive match")
	}
}

func TestEnvFilter_Match_Miss(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"production"})
	if f.Match(envEntry("env", "development")) {
		t.Error("expected no match for 'development'")
	}
}

func TestEnvFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"production"})
	if f.Match(map[string]interface{}{"other": "production"}) {
		t.Error("expected no match when field is absent")
	}
}

func TestEnvFilter_Match_NonStringValue(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"production"})
	if f.Match(map[string]interface{}{"env": 42}) {
		t.Error("expected no match for non-string value")
	}
}
