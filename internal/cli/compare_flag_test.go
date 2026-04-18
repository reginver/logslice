package cli

import (
	"testing"
)

func TestParseComparePairs_Valid(t *testing.T) {
	pairs := []string{"latency:gt:100", "code:eq:200"}
	filters, err := parseComparePairs(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseComparePairs_Empty(t *testing.T) {
	filters, err := parseComparePairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters")
	}
}

func TestParseComparePair_MissingParts(t *testing.T) {
	_, err := parseComparePair("latency:gt")
	if err == nil {
		t.Fatal("expected error for missing value part")
	}
}

func TestParseComparePair_EmptyField(t *testing.T) {
	_, err := parseComparePair(":gt:100")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseComparePair_InvalidValue(t *testing.T) {
	_, err := parseComparePair("latency:gt:notanumber")
	if err == nil {
		t.Fatal("expected error for non-numeric value")
	}
}

func TestParseComparePair_UnknownOp(t *testing.T) {
	_, err := parseComparePair("latency:bad:100")
	if err == nil {
		t.Fatal("expected error for unknown operator")
	}
}
