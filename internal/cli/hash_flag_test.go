package cli

import (
	"testing"
)

func TestParseHashPairs_Valid(t *testing.T) {
	filters, err := parseHashPairs([]string{"request_id=ab12", "user=ff"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseHashPairs_Empty(t *testing.T) {
	filters, err := parseHashPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseHashPair_MissingParts(t *testing.T) {
	_, err := parseHashPair("noequalssign")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseHashPair_EmptyField(t *testing.T) {
	_, err := parseHashPair("=ab12")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseHashPair_EmptyPrefix(t *testing.T) {
	_, err := parseHashPair("request_id=")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestParseHashPair_PrefixTooLong(t *testing.T) {
	_, err := parseHashPair("request_id=aabbccddeeff00112233445566778899x")
	if err == nil {
		t.Fatal("expected error for prefix exceeding MD5 hex length")
	}
}
