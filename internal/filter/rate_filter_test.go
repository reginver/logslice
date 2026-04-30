package filter

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func rateEntry(ts time.Time, field string, value float64) parser.LogEntry {
	return parser.LogEntry{
		Timestamp: &ts,
		Fields:    map[string]interface{}{field: value},
	}
}

func TestNewRateFilter_Valid(t *testing.T) {
	f, err := NewRateFilter("bytes", time.Second, 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewRateFilter_EmptyField(t *testing.T) {
	_, err := NewRateFilter("", time.Second, 0, 100)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRateFilter_ZeroWindow(t *testing.T) {
	_, err := NewRateFilter("bytes", 0, 0, 100)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewRateFilter_MinGreaterThanMax(t *testing.T) {
	_, err := NewRateFilter("bytes", time.Second, 200, 100)
	if err == nil {
		t.Fatal("expected error when minRate > maxRate")
	}
}

func TestRateFilter_InsufficientPoints(t *testing.T) {
	f, _ := NewRateFilter("bytes", 5*time.Second, 0, 1000)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// Only one point — cannot compute rate.
	if f.Match(rateEntry(base, "bytes", 50)) {
		t.Fatal("expected false with only one point")
	}
}

func TestRateFilter_Match_WithinBounds(t *testing.T) {
	f, _ := NewRateFilter("bytes", 10*time.Second, 5, 50)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// Two points 2 s apart, values 10 and 10 → rate = 20/2 = 10 per second.
	f.Match(rateEntry(base, "bytes", 10))
	result := f.Match(rateEntry(base.Add(2*time.Second), "bytes", 10))
	if !result {
		t.Fatal("expected match within bounds")
	}
}

func TestRateFilter_Match_BelowMin(t *testing.T) {
	f, _ := NewRateFilter("bytes", 10*time.Second, 100, 0)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// rate = 20/2 = 10 per second — below minRate 100.
	f.Match(rateEntry(base, "bytes", 10))
	result := f.Match(rateEntry(base.Add(2*time.Second), "bytes", 10))
	if result {
		t.Fatal("expected no match below minRate")
	}
}

func TestRateFilter_Match_AboveMax(t *testing.T) {
	f, _ := NewRateFilter("bytes", 10*time.Second, 0, 5)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// rate = 20/2 = 10 per second — above maxRate 5.
	f.Match(rateEntry(base, "bytes", 10))
	result := f.Match(rateEntry(base.Add(2*time.Second), "bytes", 10))
	if result {
		t.Fatal("expected no match above maxRate")
	}
}

func TestRateFilter_NoTimestamp(t *testing.T) {
	f, _ := NewRateFilter("bytes", time.Second, 0, 1000)
	e := parser.LogEntry{Fields: map[string]interface{}{"bytes": 10.0}}
	if f.Match(e) {
		t.Fatal("expected false when timestamp is nil")
	}
}

func TestRateFilter_MissingField(t *testing.T) {
	f, _ := NewRateFilter("bytes", time.Second, 0, 1000)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e := parser.LogEntry{Timestamp: &base, Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Fatal("expected false when field is absent")
	}
}
