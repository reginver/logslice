package filter

import (
	"testing"

	"github.com/example/logslice/internal/parser"
)

func numEntry(field string, val interface{}) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]interface{}{field: val}}
}

func TestNewRangeFilter_Valid(t *testing.T) {
	f, err := NewRangeFilter("latency", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "latency" {
		t.Errorf("expected field latency, got %s", f.Field)
	}
}

func TestNewRangeFilter_EmptyField(t *testing.T) {
	_, err := NewRangeFilter("", 0, 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRangeFilter_MinGreaterThanMax(t *testing.T) {
	_, err := NewRangeFilter("size", 50, 10)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestRangeFilter_Match_Float64(t *testing.T) {
	f, _ := NewRangeFilter("latency", 10, 200)
	if !f.Match(numEntry("latency", float64(50))) {
		t.Error("expected match for value 50")
	}
	if f.Match(numEntry("latency", float64(5))) {
		t.Error("expected no match for value 5")
	}
}

func TestRangeFilter_Match_Int(t *testing.T) {
	f, _ := NewRangeFilter("code", 200, 299)
	if !f.Match(numEntry("code", 200)) {
		t.Error("expected match for int 200")
	}
	if f.Match(numEntry("code", 500)) {
		t.Error("expected no match for int 500")
	}
}

func TestRangeFilter_Match_MissingField(t *testing.T) {
	f, _ := NewRangeFilter("latency", 0, 100)
	e := parser.LogEntry{Fields: map[string]interface{}{"other": float64(50)}}
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestRangeFilter_Match_NonNumeric(t *testing.T) {
	f, _ := NewRangeFilter("latency", 0, 100)
	if f.Match(numEntry("latency", "fast")) {
		t.Error("expected no match for non-numeric value")
	}
}

func TestRangeFilter_BoundaryInclusive(t *testing.T) {
	f, _ := NewRangeFilter("score", 1, 5)
	if !f.Match(numEntry("score", float64(1))) {
		t.Error("expected match at lower bound")
	}
	if !f.Match(numEntry("score", float64(5))) {
		t.Error("expected match at upper bound")
	}
}
