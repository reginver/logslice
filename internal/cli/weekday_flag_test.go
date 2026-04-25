package cli

import (
	"testing"
)

func TestParseWeekdayPairs_Valid(t *testing.T) {
	filters, err := parseWeekdayPairs([]string{"ts=monday,friday", "created=saturday"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseWeekdayPairs_Empty(t *testing.T) {
	filters, err := parseWeekdayPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseWeekdayPair_MissingEquals(t *testing.T) {
	_, err := parseWeekdayPair("ts-monday")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseWeekdayPair_EmptyField(t *testing.T) {
	_, err := parseWeekdayPair("=monday")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseWeekdayPair_EmptyDays(t *testing.T) {
	_, err := parseWeekdayPair("ts=")
	if err == nil {
		t.Fatal("expected error for empty days list")
	}
}

func TestParseWeekdayPair_UnknownDay(t *testing.T) {
	_, err := parseWeekdayPair("ts=funday")
	if err == nil {
		t.Fatal("expected error for unknown weekday")
	}
}

func TestParseWeekdayPair_CaseInsensitive(t *testing.T) {
	_, err := parseWeekdayPair("ts=Monday,FRIDAY")
	if err != nil {
		t.Fatalf("expected case-insensitive match, got error: %v", err)
	}
}
