package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func containsEntry(field, value string) *parser.Entry {
	return &parser.Entry{Fields: map[string]interface{}{field: value}}
}

func TestNewContainsFilter_Valid(t *testing.T) {
	_, err := filter.NewContainsFilter("msg", []string{"error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewContainsFilter_EmptyField(t *testing.T) {
	_, err := filter.NewContainsFilter("", []string{"error"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewContainsFilter_NoValues(t *testing.T) {
	_, err := filter.NewContainsFilter("msg", []string{})
	if err == nil {
		t.Fatal("expected error for empty values slice")
	}
}

func TestNewContainsFilter_EmptyValue(t *testing.T) {
	_, err := filter.NewContainsFilter("msg", []string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for empty value element")
	}
}

func TestContainsFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", []string{"error"})
	e := containsEntry("msg", "disk error occurred")
	if !f.Match(e) {
		t.Fatal("expected match")
	}
}

func TestContainsFilter_Match_Miss(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", []string{"error"})
	e := containsEntry("msg", "all systems nominal")
	if f.Match(e) {
		t.Fatal("expected no match")
	}
}

func TestContainsFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", []string{"ERROR"})
	e := containsEntry("msg", "disk error occurred")
	if !f.Match(e) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestContainsFilter_Match_MultipleValues(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", []string{"disk", "error"})
	e := containsEntry("msg", "disk error occurred")
	if !f.Match(e) {
		t.Fatal("expected match for all values present")
	}
	e2 := containsEntry("msg", "disk warning")
	if f.Match(e2) {
		t.Fatal("expected no match when only one value present")
	}
}

func TestContainsFilter_Match_MissingField(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", []string{"error"})
	e := &parser.Entry{Fields: map[string]interface{}{"level": "info"}}
	if f.Match(e) {
		t.Fatal("expected no match for missing field")
	}
}
