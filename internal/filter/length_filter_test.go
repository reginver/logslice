package filter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/parser"
)

func lengthEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]interface{}{field: value}}
}

func TestNewLengthFilter_Valid(t *testing.T) {
	_, err := filter.NewLengthFilter("msg", 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLengthFilter_EmptyField(t *testing.T) {
	_, err := filter.NewLengthFilter("", 1, 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewLengthFilter_MinGreaterThanMax(t *testing.T) {
	_, err := filter.NewLengthFilter("msg", 10, 5)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNewLengthFilter_UnboundedMin(t *testing.T) {
	_, err := filter.NewLengthFilter("msg", -1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLengthFilter_UnboundedMax(t *testing.T) {
	_, err := filter.NewLengthFilter("msg", 2, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLengthFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewLengthFilter("msg", 3, 10)
	e := lengthEntry("msg", "hello")
	if !f.Match(e) {
		t.Fatal("expected match")
	}
}

func TestLengthFilter_Match_TooShort(t *testing.T) {
	f, _ := filter.NewLengthFilter("msg", 10, 20)
	e := lengthEntry("msg", "hi")
	if f.Match(e) {
		t.Fatal("expected no match")
	}
}

func TestLengthFilter_Match_TooLong(t *testing.T) {
	f, _ := filter.NewLengthFilter("msg", 1, 3)
	e := lengthEntry("msg", "toolong")
	if f.Match(e) {
		t.Fatal("expected no match")
	}
}

func TestLengthFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewLengthFilter("msg", 1, 10)
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Fatal("expected no match for absent field")
	}
}

func TestLengthFilter_Match_Unbounded(t *testing.T) {
	f, _ := filter.NewLengthFilter("msg", -1, -1)
	e := lengthEntry("msg", "anything at all")
	if !f.Match(e) {
		t.Fatal("expected match with no bounds")
	}
}
