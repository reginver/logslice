package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func latencyEntry(field string, ms float64) parser.Entry {
	return parser.Entry{Fields: map[string]any{field: ms}}
}

func TestNewLatencyFilter_Valid(t *testing.T) {
	_, err := filter.NewLatencyFilter("latency_ms", 0, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLatencyFilter_EmptyField(t *testing.T) {
	_, err := filter.NewLatencyFilter("", 0, 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewLatencyFilter_BothZero(t *testing.T) {
	_, err := filter.NewLatencyFilter("latency_ms", 0, 0)
	if err == nil {
		t.Fatal("expected error when both bounds are zero")
	}
}

func TestNewLatencyFilter_MinGreaterThanMax(t *testing.T) {
	_, err := filter.NewLatencyFilter("latency_ms", 200*time.Millisecond, 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestLatencyFilter_Match_WithinRange(t *testing.T) {
	f, _ := filter.NewLatencyFilter("latency_ms", 100*time.Millisecond, 500*time.Millisecond)
	if !f.Match(latencyEntry("latency_ms", 250)) {
		t.Error("expected match for value within range")
	}
}

func TestLatencyFilter_Match_BelowMin(t *testing.T) {
	f, _ := filter.NewLatencyFilter("latency_ms", 100*time.Millisecond, 500*time.Millisecond)
	if f.Match(latencyEntry("latency_ms", 50)) {
		t.Error("expected no match for value below min")
	}
}

func TestLatencyFilter_Match_AboveMax(t *testing.T) {
	f, _ := filter.NewLatencyFilter("latency_ms", 100*time.Millisecond, 500*time.Millisecond)
	if f.Match(latencyEntry("latency_ms", 600)) {
		t.Error("expected no match for value above max")
	}
}

func TestLatencyFilter_Match_NoUpperBound(t *testing.T) {
	f, _ := filter.NewLatencyFilter("latency_ms", 100*time.Millisecond, 0)
	if !f.Match(latencyEntry("latency_ms", 99999)) {
		t.Error("expected match with no upper bound")
	}
	if f.Match(latencyEntry("latency_ms", 10)) {
		t.Error("expected no match below min")
	}
}

func TestLatencyFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewLatencyFilter("latency_ms", 0, 500*time.Millisecond)
	e := parser.Entry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestLatencyFilter_Match_NonNumericField(t *testing.T) {
	f, _ := filter.NewLatencyFilter("latency_ms", 0, 500*time.Millisecond)
	e := parser.Entry{Fields: map[string]any{"latency_ms": "fast"}}
	if f.Match(e) {
		t.Error("expected no match for non-numeric field value")
	}
}
