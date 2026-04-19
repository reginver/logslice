package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func boolEntry(field string, value interface{}) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]interface{}{field: value}}
}

func TestNewBoolFilter_Valid(t *testing.T) {
	_, err := filter.NewBoolFilter("active", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewBoolFilter_EmptyField(t *testing.T) {
	_, err := filter.NewBoolFilter("", true)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestBoolFilter_Match_BoolTrue(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", true)
	if !f.Match(boolEntry("active", true)) {
		t.Error("expected match")
	}
}

func TestBoolFilter_Match_BoolFalse(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", false)
	if !f.Match(boolEntry("active", false)) {
		t.Error("expected match")
	}
}

func TestBoolFilter_Match_StringTrue(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", true)
	for _, v := range []string{"true", "1", "yes"} {
		if !f.Match(boolEntry("active", v)) {
			t.Errorf("expected match for string %q", v)
		}
	}
}

func TestBoolFilter_Match_StringFalse(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", false)
	for _, v := range []string{"false", "0", "no"} {
		if !f.Match(boolEntry("active", v)) {
			t.Errorf("expected match for string %q", v)
		}
	}
}

func TestBoolFilter_Match_Float64(t *testing.T) {
	f, _ := filter.NewBoolFilter("flag", true)
	if !f.Match(boolEntry("flag", float64(1))) {
		t.Error("expected match for non-zero float")
	}
	f2, _ := filter.NewBoolFilter("flag", false)
	if !f2.Match(boolEntry("flag", float64(0))) {
		t.Error("expected match for zero float")
	}
}

func TestBoolFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", true)
	if f.Match(parser.LogEntry{Fields: map[string]interface{}{}}) {
		t.Error("expected no match for absent field")
	}
}
