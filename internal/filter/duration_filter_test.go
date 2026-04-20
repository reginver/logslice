package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func durationEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewDurationFilter_Valid(t *testing.T) {
	_, err := filter.NewDurationFilter("latency", 0, 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDurationFilter_EmptyField(t *testing.T) {
	_, err := filter.NewDurationFilter("", 0, time.Second)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewDurationFilter_BothZero(t *testing.T) {
	_, err := filter.NewDurationFilter("latency", 0, 0)
	if err == nil {
		t.Fatal("expected error when both bounds are zero")
	}
}

func TestNewDurationFilter_MinGreaterThanMax(t *testing.T) {
	_, err := filter.NewDurationFilter("latency", 5*time.Second, time.Second)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestDurationFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", 100*time.Millisecond, 2*time.Second)
	if !f.Match(durationEntry("latency", "500ms")) {
		t.Error("expected match for 500ms within [100ms, 2s]")
	}
}

func TestDurationFilter_Match_BelowMin(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", 500*time.Millisecond, 2*time.Second)
	if f.Match(durationEntry("latency", "100ms")) {
		t.Error("expected no match for 100ms below min 500ms")
	}
}

func TestDurationFilter_Match_AboveMax(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", 0, time.Second)
	if f.Match(durationEntry("latency", "5s")) {
		t.Error("expected no match for 5s above max 1s")
	}
}

func TestDurationFilter_Match_InvalidValue(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", 0, time.Second)
	if f.Match(durationEntry("latency", "not-a-duration")) {
		t.Error("expected no match for unparseable value")
	}
}

func TestDurationFilter_Match_MissingField(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", 0, time.Second)
	entry := parser.LogEntry{Fields: map[string]any{"other": "1s"}}
	if f.Match(entry) {
		t.Error("expected no match when field is absent")
	}
}

func TestDurationFilter_Match_OnlyMin(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", time.Second, 0)
	if !f.Match(durationEntry("latency", "10s")) {
		t.Error("expected match for 10s with only min=1s set")
	}
	if f.Match(durationEntry("latency", "100ms")) {
		t.Error("expected no match for 100ms below min=1s")
	}
}
