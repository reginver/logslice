package cli

import (
	"testing"
)

func TestParseDatePairs_Valid(t *testing.T) {
	filters, err := parseDatePairs([]string{"2024-03-15", "time=2024-06-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseDatePairs_Empty(t *testing.T) {
	filters, err := parseDatePairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseDatePair_BareDate(t *testing.T) {
	f, err := parseDatePair("2024-01-20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestParseDatePair_EmptyValue(t *testing.T) {
	_, err := parseDatePair("")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestParseDatePair_EmptyField(t *testing.T) {
	_, err := parseDatePair("=2024-03-15")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseDatePair_InvalidDate(t *testing.T) {
	_, err := parseDatePair("not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestParseDatePair_EmptyDate(t *testing.T) {
	_, err := parseDatePair("time=")
	if err == nil {
		t.Fatal("expected error for empty date value")
	}
}
