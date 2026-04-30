package cli

import (
	"testing"
)

func TestParsePercentilePairs_Valid(t *testing.T) {
	filters, err := parsePercentilePairs([]string{"latency=10-90", "size=0-50"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParsePercentilePairs_Empty(t *testing.T) {
	filters, err := parsePercentilePairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParsePercentilePair_MissingEquals(t *testing.T) {
	_, err := parsePercentilePair("latency10-90")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParsePercentilePair_EmptyField(t *testing.T) {
	_, err := parsePercentilePair("=10-90")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParsePercentilePair_MissingDash(t *testing.T) {
	_, err := parsePercentilePair("latency=1090")
	if err == nil {
		t.Fatal("expected error for missing '-' in range")
	}
}

func TestParsePercentilePair_InvalidMin(t *testing.T) {
	_, err := parsePercentilePair("latency=abc-90")
	if err == nil {
		t.Fatal("expected error for non-numeric min")
	}
}

func TestParsePercentilePair_InvalidMax(t *testing.T) {
	_, err := parsePercentilePair("latency=10-xyz")
	if err == nil {
		t.Fatal("expected error for non-numeric max")
	}
}
