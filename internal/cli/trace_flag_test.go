package cli

import (
	"testing"
)

func TestParseTracePairs_Valid(t *testing.T) {
	filters, err := parseTracePairs([]string{"trace_id=abc123,def*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(filters))
	}
}

func TestParseTracePairs_Empty(t *testing.T) {
	filters, err := parseTracePairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseTracePair_MissingEquals(t *testing.T) {
	_, err := parseTracePair("trace_idabc123")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseTracePair_EmptyField(t *testing.T) {
	_, err := parseTracePair("=abc123")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseTracePair_EmptySpec(t *testing.T) {
	_, err := parseTracePair("trace_id=")
	if err == nil {
		t.Fatal("expected error for empty spec")
	}
}

func TestParseTracePair_PrefixPattern(t *testing.T) {
	f, err := parseTracePair("trace_id=span-*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestParseTracePairs_Multiple(t *testing.T) {
	filters, err := parseTracePairs([]string{
		"trace_id=abc123",
		"span_id=sp-*,sp-002",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}
