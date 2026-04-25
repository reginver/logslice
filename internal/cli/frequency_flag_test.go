package cli

import (
	"testing"
)

func TestParseFrequencyPairs_Valid(t *testing.T) {
	filters, err := parseFrequencyPairs([]string{"host:2:10", "svc:1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseFrequencyPairs_Empty(t *testing.T) {
	filters, err := parseFrequencyPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseFrequencyPair_MissingParts(t *testing.T) {
	_, err := parseFrequencyPair("host")
	if err == nil {
		t.Fatal("expected error for missing min count")
	}
}

func TestParseFrequencyPair_EmptyField(t *testing.T) {
	_, err := parseFrequencyPair(":2:5")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseFrequencyPair_InvalidMin(t *testing.T) {
	_, err := parseFrequencyPair("host:abc")
	if err == nil {
		t.Fatal("expected error for non-numeric min")
	}
}

func TestParseFrequencyPair_InvalidMax(t *testing.T) {
	_, err := parseFrequencyPair("host:2:xyz")
	if err == nil {
		t.Fatal("expected error for non-numeric max")
	}
}

func TestParseFrequencyPair_UnboundedMax(t *testing.T) {
	f, err := parseFrequencyPair("host:3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
