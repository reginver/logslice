package filter

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func todEntry(ts time.Time) *parser.Entry {
	return &parser.Entry{Timestamp: &ts}
}

func TestNewTimeOfDayFilter_Valid(t *testing.T) {
	f, err := NewTimeOfDayFilter("09:00", "17:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewTimeOfDayFilter_EmptyBounds(t *testing.T) {
	if _, err := NewTimeOfDayFilter("", "17:00"); err == nil {
		t.Fatal("expected error for empty start")
	}
	if _, err := NewTimeOfDayFilter("09:00", ""); err == nil {
		t.Fatal("expected error for empty end")
	}
}

func TestNewTimeOfDayFilter_StartAfterEnd(t *testing.T) {
	if _, err := NewTimeOfDayFilter("17:00", "09:00"); err == nil {
		t.Fatal("expected error when start > end")
	}
}

func TestNewTimeOfDayFilter_InvalidFormat(t *testing.T) {
	if _, err := NewTimeOfDayFilter("9am", "17:00"); err == nil {
		t.Fatal("expected error for bad format")
	}
}

func TestTimeOfDayFilter_Match_Inside(t *testing.T) {
	f, _ := NewTimeOfDayFilter("09:00", "17:00")
	ts := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	if !f.Match(todEntry(ts)) {
		t.Error("expected match for time inside window")
	}
}

func TestTimeOfDayFilter_Match_Outside(t *testing.T) {
	f, _ := NewTimeOfDayFilter("09:00", "17:00")
	ts := time.Date(2024, 1, 15, 20, 0, 0, 0, time.UTC)
	if f.Match(todEntry(ts)) {
		t.Error("expected no match for time outside window")
	}
}

func TestTimeOfDayFilter_Match_Boundary(t *testing.T) {
	f, _ := NewTimeOfDayFilter("09:00", "17:00")
	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)
	if !f.Match(todEntry(start)) {
		t.Error("expected match on start boundary")
	}
	if !f.Match(todEntry(end)) {
		t.Error("expected match on end boundary")
	}
}

func TestTimeOfDayFilter_Match_NilTimestamp(t *testing.T) {
	f, _ := NewTimeOfDayFilter("09:00", "17:00")
	if f.Match(&parser.Entry{}) {
		t.Error("expected no match for nil timestamp")
	}
}

func TestTimeOfDayFilter_WithSeconds(t *testing.T) {
	f, err := NewTimeOfDayFilter("09:00:00", "09:00:30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hit := time.Date(2024, 1, 15, 9, 0, 15, 0, time.UTC)
	miss := time.Date(2024, 1, 15, 9, 0, 45, 0, time.UTC)
	if !f.Match(todEntry(hit)) {
		t.Error("expected match within second-precision window")
	}
	if f.Match(todEntry(miss)) {
		t.Error("expected no match outside second-precision window")
	}
}
