package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func suffixEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewSuffixFilter_Valid(t *testing.T) {
	f, err := NewSuffixFilter("msg", ".go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewSuffixFilter_EmptyField(t *testing.T) {
	_, err := NewSuffixFilter("", ".go")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSuffixFilter_EmptySuffix(t *testing.T) {
	_, err := NewSuffixFilter("msg", "")
	if err == nil {
		t.Fatal("expected error for empty suffix")
	}
}

func TestSuffixFilter_Match_Hit(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".go")
	e := suffixEntry("file", "main.go")
	if !f.Match(e) {
		t.Error("expected match")
	}
}

func TestSuffixFilter_Match_Miss(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".go")
	e := suffixEntry("file", "main.txt")
	if f.Match(e) {
		t.Error("expected no match")
	}
}

func TestSuffixFilter_Match_MissingField(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".go")
	e := parser.LogEntry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Error("expected no match for missing field")
	}
}

func TestSuffixFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewSuffixFilter("code", "404")
	e := parser.LogEntry{Fields: map[string]any{"code": 404}}
	if !f.Match(e) {
		t.Error("expected match via fmt fallback")
	}
}
