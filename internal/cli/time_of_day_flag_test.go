package cli

import (
	"testing"
)

func TestParseTimeOfDayPairs_Valid(t *testing.T) {
	filters, err := parseTimeOfDayPairs([]string{"09:00-17:00", "22:00-23:59"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseTimeOfDayPairs_Empty(t *testing.T) {
	filters, err := parseTimeOfDayPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseTimeOfDayPair_MissingDash(t *testing.T) {
	if _, err := parseTimeOfDayPair("0900"); err == nil {
		t.Fatal("expected error for missing dash separator")
	}
}

func TestParseTimeOfDayPair_EmptyStart(t *testing.T) {
	if _, err := parseTimeOfDayPair("-17:00"); err == nil {
		t.Fatal("expected error for empty start")
	}
}

func TestParseTimeOfDayPair_EmptyEnd(t *testing.T) {
	if _, err := parseTimeOfDayPair("09:00-"); err == nil {
		t.Fatal("expected error for empty end")
	}
}

func TestParseTimeOfDayPair_InvalidFormat(t *testing.T) {
	if _, err := parseTimeOfDayPair("9am-5pm"); err == nil {
		t.Fatal("expected error for non-HH:MM format")
	}
}

func TestParseTimeOfDayPair_WithSeconds(t *testing.T) {
	f, err := parseTimeOfDayPair("08:30:00-08:30:59")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
