package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func keywordEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewKeywordFilter_Valid(t *testing.T) {
	f, err := NewKeywordFilter("msg", []string{"error", "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewKeywordFilter_EmptyField(t *testing.T) {
	_, err := NewKeywordFilter("", []string{"error"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewKeywordFilter_NoKeywords(t *testing.T) {
	_, err := NewKeywordFilter("msg", []string{})
	if err == nil {
		t.Fatal("expected error for empty keyword list")
	}
}

func TestNewKeywordFilter_EmptyKeyword(t *testing.T) {
	_, err := NewKeywordFilter("msg", []string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for empty keyword")
	}
}

func TestKeywordFilter_Match_Hit(t *testing.T) {
	f, _ := NewKeywordFilter("msg", []string{"error"})
	e := keywordEntry("msg", "connection error occurred")
	if !f.Match(e) {
		t.Error("expected match")
	}
}

func TestKeywordFilter_Match_Miss(t *testing.T) {
	f, _ := NewKeywordFilter("msg", []string{"error"})
	e := keywordEntry("msg", "all systems nominal")
	if f.Match(e) {
		t.Error("expected no match")
	}
}

func TestKeywordFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := NewKeywordFilter("msg", []string{"error"})
	e := keywordEntry("msg", "CRITICAL ERROR")
	if !f.Match(e) {
		t.Error("expected case-insensitive match")
	}
}

func TestKeywordFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewKeywordFilter("msg", []string{"error"})
	e := parser.LogEntry{Fields: map[string]any{"level": "info"}}
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestKeywordFilter_Match_MultipleKeywords(t *testing.T) {
	f, _ := NewKeywordFilter("msg", []string{"timeout", "refused"})
	if !f.Match(keywordEntry("msg", "connection refused")) {
		t.Error("expected match on second keyword")
	}
	if f.Match(keywordEntry("msg", "all good")) {
		t.Error("expected no match")
	}
}
