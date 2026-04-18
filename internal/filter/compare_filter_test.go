package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func cmpEntry(field string, val interface{}) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]interface{}{field: val}}
}

func TestNewCompareFilter_Valid(t *testing.T) {
	_, err := NewCompareFilter("latency", OpGt, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCompareFilter_EmptyField(t *testing.T) {
	_, err := NewCompareFilter("", OpEq, 0)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewCompareFilter_UnknownOp(t *testing.T) {
	_, err := NewCompareFilter("x", CompareOp("bad"), 0)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestCompareFilter_Eq(t *testing.T) {
	f, _ := NewCompareFilter("code", OpEq, 200)
	if !f.Match(cmpEntry("code", float64(200))) {
		t.Error("expected match")
	}
	if f.Match(cmpEntry("code", float64(404))) {
		t.Error("expected no match")
	}
}

func TestCompareFilter_Gt(t *testing.T) {
	f, _ := NewCompareFilter("latency", OpGt, 100)
	if !f.Match(cmpEntry("latency", float64(200))) {
		t.Error("expected match for 200 > 100")
	}
	if f.Match(cmpEntry("latency", float64(50))) {
		t.Error("expected no match for 50 > 100")
	}
}

func TestCompareFilter_StringNumeric(t *testing.T) {
	f, _ := NewCompareFilter("score", OpLte, 5.0)
	if !f.Match(cmpEntry("score", "3.5")) {
		t.Error("expected match for string '3.5' <= 5.0")
	}
}

func TestCompareFilter_MissingField(t *testing.T) {
	f, _ := NewCompareFilter("x", OpEq, 1)
	if f.Match(parser.LogEntry{Fields: map[string]interface{}{}}) {
		t.Error("expected no match for missing field")
	}
}

func TestCompareFilter_NonNumericString(t *testing.T) {
	f, _ := NewCompareFilter("val", OpGt, 0)
	if f.Match(cmpEntry("val", "notanumber")) {
		t.Error("expected no match for non-numeric string")
	}
}
