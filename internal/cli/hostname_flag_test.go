package cli

import "testing"

func TestParseHostnamePairs_Valid(t *testing.T) {
	filters, err := parseHostnamePairs([]string{"host=web-*,db-*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(filters))
	}
}

func TestParseHostnamePairs_Empty(t *testing.T) {
	filters, err := parseHostnamePairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseHostnamePair_MissingEquals(t *testing.T) {
	_, err := parseHostnamePair("hostweb-01")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseHostnamePair_EmptyField(t *testing.T) {
	_, err := parseHostnamePair("=web-*")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseHostnamePair_EmptyPatternList(t *testing.T) {
	_, err := parseHostnamePair("host=")
	if err == nil {
		t.Fatal("expected error for empty pattern list")
	}
}

func TestParseHostnamePair_MultiplePatterns(t *testing.T) {
	f, err := parseHostnamePair("host=web-*,db-01,cache-*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
