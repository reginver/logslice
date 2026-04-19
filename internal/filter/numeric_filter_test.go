package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func numericEntry(field string, value interface{}) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]interface{}{field: value}}
}

func TestNewNumericFilter_Valid(t *testing.T) {
	_, err := NewNumericFilter("code", 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewNumericFilter_EmptyField(t *testing.T) {
	_, err := NewNumericFilter("", 2, 0)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewNumericFilter_ZeroDivisor(t *testing.T) {
	_, err := NewNumericFilter("code", 0, 0)
	if err == nil {
		t.Fatal("expected error for zero divisor")
	}
}

func TestNewNumericFilter_InvalidRemainder(t *testing.T) {
	_, err := NewNumericFilter("code", 3, 5)
	if err == nil {
		t.Fatal("expected error when remainder >= divisor")
	}
}

func TestNumericFilter_Match_Float64Even(t *testing.T) {
	f, _ := NewNumericFilter("val", 2, 0)
	if !f.Match(numericEntry("val", float64(4))) {
		t.Error("expected match for even number")
	}
	if f.Match(numericEntry("val", float64(3))) {
		t.Error("expected no match for odd number")
	}
}

func TestNumericFilter_Match_StringValue(t *testing.T) {
	f, _ := NewNumericFilter("val", 3, 1)
	if !f.Match(numericEntry("val", "7")) {
		t.Error("expected match: 7 mod 3 == 1")
	}
	if f.Match(numericEntry("val", "6")) {
		t.Error("expected no match: 6 mod 3 == 0")
	}
}

func TestNumericFilter_Match_MissingField(t *testing.T) {
	f, _ := NewNumericFilter("val", 2, 0)
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Error("expected no match for missing field")
	}
}

func TestNumericFilter_Match_NonNumericString(t *testing.T) {
	f, _ := NewNumericFilter("val", 2, 0)
	if f.Match(numericEntry("val", "abc")) {
		t.Error("expected no match for non-numeric string")
	}
}
