package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func subEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewSubstringFilter_Valid(t *testing.T) {
	f, err := NewSubstringFilter("msg", "error", false)
	if err != nil || f == nil {
		t.Fatalf("expected valid filter, got err=%v", err)
	}
}

func TestNewSubstringFilter_EmptyField(t *testing.T) {
	_, err := NewSubstringFilter("", "error", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSubstringFilter_EmptySubstring(t *testing.T) {
	_, err := NewSubstringFilter("msg", "", false)
	if err == nil {
		t.Fatal("expected error for empty substring")
	}
}

func TestSubstringFilter_Match_Hit(t *testing.T) {
	f, _ := NewSubstringFilter("msg", "error", false)
	if !f.Match(subEntry("msg", "an error occurred")) {
		t.Error("expected match")
	}
}

func TestSubstringFilter_Match_Miss(t *testing.T) {
	f, _ := NewSubstringFilter("msg", "error", false)
	if f.Match(subEntry("msg", "all good")) {
		t.Error("expected no match")
	}
}

func TestSubstringFilter_Match_CaseFold(t *testing.T) {
	f, _ := NewSubstringFilter("msg", "ERROR", true)
	if !f.Match(subEntry("msg", "an error occurred")) {
		t.Error("expected case-insensitive match")
	}
}

func TestSubstringFilter_Match_MissingField(t *testing.T) {
	f, _ := NewSubstringFilter("msg", "error", false)
	e := parser.LogEntry{Fields: map[string]any{"other": "error"}}
	if f.Match(e) {
		t.Error("expected no match for missing field")
	}
}

func TestSubstringFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewSubstringFilter("code", "404", false)
	e := parser.LogEntry{Fields: map[string]any{"code": 404}}
	if !f.Match(e) {
		t.Error("expected match on stringified numeric field")
	}
}
