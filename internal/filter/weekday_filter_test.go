package filter

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func weekdayEntry(field string, t time.Time) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]interface{}{field: t}}
}

func TestNewWeekdayFilter_Valid(t *testing.T) {
	f, err := NewWeekdayFilter("ts", []string{"Monday", "friday"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewWeekdayFilter_EmptyField(t *testing.T) {
	_, err := NewWeekdayFilter("", []string{"monday"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewWeekdayFilter_NoDays(t *testing.T) {
	_, err := NewWeekdayFilter("ts", []string{})
	if err == nil {
		t.Fatal("expected error for empty days list")
	}
}

func TestNewWeekdayFilter_UnknownDay(t *testing.T) {
	_, err := NewWeekdayFilter("ts", []string{"funday"})
	if err == nil {
		t.Fatal("expected error for unknown weekday")
	}
}

func TestWeekdayFilter_Match_Hit(t *testing.T) {
	// 2024-01-01 is a Monday
	monday := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	f, _ := NewWeekdayFilter("ts", []string{"monday"})
	if !f.Match(weekdayEntry("ts", monday)) {
		t.Error("expected match on Monday")
	}
}

func TestWeekdayFilter_Match_Miss(t *testing.T) {
	// 2024-01-02 is a Tuesday
	tuesday := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)
	f, _ := NewWeekdayFilter("ts", []string{"monday"})
	if f.Match(weekdayEntry("ts", tuesday)) {
		t.Error("expected no match on Tuesday when only Monday allowed")
	}
}

func TestWeekdayFilter_Match_StringField(t *testing.T) {
	f, _ := NewWeekdayFilter("ts", []string{"monday"})
	e := parser.LogEntry{Fields: map[string]interface{}{"ts": "2024-01-01T10:00:00Z"}}
	if !f.Match(e) {
		t.Error("expected match for RFC3339 string on Monday")
	}
}

func TestWeekdayFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewWeekdayFilter("ts", []string{"monday"})
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestWeekdayFilter_Match_MultipleAllowed(t *testing.T) {
	f, _ := NewWeekdayFilter("ts", []string{"saturday", "sunday"})
	sat := time.Date(2024, 1, 6, 9, 0, 0, 0, time.UTC) // Saturday
	sun := time.Date(2024, 1, 7, 9, 0, 0, 0, time.UTC) // Sunday
	if !f.Match(weekdayEntry("ts", sat)) {
		t.Error("expected match on Saturday")
	}
	if !f.Match(weekdayEntry("ts", sun)) {
		t.Error("expected match on Sunday")
	}
}
