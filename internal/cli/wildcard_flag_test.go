package cli

import (
	"testing"
)

func TestParseWildcardPairs_Valid(t *testing.T) {
	pairs, err := parseWildcardPairs([]string{"service=auth-*", "host=web-?"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].field != "service" || pairs[0].pattern != "auth-*" {
		t.Errorf("unexpected pair[0]: %+v", pairs[0])
	}
}

func TestParseWildcardPairs_Empty(t *testing.T) {
	pairs, err := parseWildcardPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestParseWildcardPair_MissingParts(t *testing.T) {
	_, err := parseWildcardPair("service")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseWildcardPair_EmptyField(t *testing.T) {
	_, err := parseWildcardPair("=auth-*")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseWildcardPair_EmptyPattern(t *testing.T) {
	_, err := parseWildcardPair("service=")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}
