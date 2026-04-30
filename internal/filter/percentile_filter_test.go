package filter

import (
	"testing"
)

func pctEntry(field string, val float64) map[string]interface{} {
	return map[string]interface{}{field: val}
}

func TestNewPercentileFilter_Valid(t *testing.T) {
	f, err := NewPercentileFilter("latency", 0, 90)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewPercentileFilter_EmptyField(t *testing.T) {
	_, err := NewPercentileFilter("", 10, 90)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewPercentileFilter_MinEqualsMax(t *testing.T) {
	_, err := NewPercentileFilter("latency", 50, 50)
	if err == nil {
		t.Fatal("expected error when minPct == maxPct")
	}
}

func TestNewPercentileFilter_MinAboveMax(t *testing.T) {
	_, err := NewPercentileFilter("latency", 80, 20)
	if err == nil {
		t.Fatal("expected error when minPct > maxPct")
	}
}

func TestNewPercentileFilter_OutOfRange(t *testing.T) {
	_, err := NewPercentileFilter("latency", -1, 101)
	if err == nil {
		t.Fatal("expected error for out-of-range percentiles")
	}
}

func TestPercentileFilter_Match_IncludesLowValues(t *testing.T) {
	f, _ := NewPercentileFilter("v", 0, 50)
	// Feed 10 values; the lower half should match.
	matched := 0
	for i := 1; i <= 10; i++ {
		if f.Match(pctEntry("v", float64(i))) {
			matched++
		}
	}
	if matched == 0 {
		t.Error("expected some low-percentile entries to match")
	}
}

func TestPercentileFilter_Match_MissingField(t *testing.T) {
	f, _ := NewPercentileFilter("latency", 0, 100)
	if f.Match(map[string]interface{}{"other": 42.0}) {
		t.Error("expected false for missing field")
	}
}

func TestPercentileFilter_Match_NonNumeric(t *testing.T) {
	f, _ := NewPercentileFilter("latency", 0, 100)
	if f.Match(map[string]interface{}{"latency": "fast"}) {
		t.Error("expected false for non-numeric value")
	}
}
