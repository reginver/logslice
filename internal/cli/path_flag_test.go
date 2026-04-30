package cli

import (
	"testing"
)

func TestParsePathPairs_Valid(t *testing.T) {
	pairs := []string{"url=/api/*", "path=/static/**"}
	filters, err := parsePathPairs(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParsePathPairs_Empty(t *testing.T) {
	filters, err := parsePathPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParsePathPair_MissingEquals(t *testing.T) {
	_, err := parsePathPair("urlapi")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParsePathPair_EmptyField(t *testing.T) {
	_, err := parsePathPair("=/api/*")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParsePathPair_EmptyPattern(t *testing.T) {
	_, err := parsePathPair("url=")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestParsePathPair_InvalidGlob(t *testing.T) {
	_, err := parsePathPair("url=[bad")
	if err == nil {
		t.Fatal("expected error for invalid glob pattern")
	}
}
