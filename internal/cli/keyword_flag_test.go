package cli

import (
	"testing"
)

func TestParseKeywordPairs_Valid(t *testing.T) {
	filters, err := parseKeywordPairs([]string{"msg=error,warn", "level=crit"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseKeywordPairs_Empty(t *testing.T) {
	filters, err := parseKeywordPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseKeywordPair_MissingEquals(t *testing.T) {
	_, err := parseKeywordPair("msgonly")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseKeywordPair_EmptyField(t *testing.T) {
	_, err := parseKeywordPair("=error")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseKeywordPair_EmptyKeyword(t *testing.T) {
	_, err := parseKeywordPair("msg=error,,warn")
	if err == nil {
		t.Fatal("expected error for empty keyword in list")
	}
}

func TestParseKeywordPair_SingleKeyword(t *testing.T) {
	f, err := parseKeywordPair("msg=timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
