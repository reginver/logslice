package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func tagEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewTagFilter_Valid(t *testing.T) {
	f, err := NewTagFilter("tags", []string{"error", "critical"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewTagFilter_EmptyField(t *testing.T) {
	_, err := NewTagFilter("", []string{"error"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTagFilter_NoTags(t *testing.T) {
	_, err := NewTagFilter("tags", []string{})
	if err == nil {
		t.Fatal("expected error for empty tag list")
	}
}

func TestNewTagFilter_EmptyTag(t *testing.T) {
	_, err := NewTagFilter("tags", []string{"error", ""})
	if err == nil {
		t.Fatal("expected error for empty tag value")
	}
}

func TestTagFilter_Match_AllPresent(t *testing.T) {
	f, _ := NewTagFilter("tags", []string{"error", "critical"})
	e := tagEntry("tags", "info, error, critical")
	if !f.Match(e) {
		t.Error("expected match when all tags present")
	}
}

func TestTagFilter_Match_MissingTag(t *testing.T) {
	f, _ := NewTagFilter("tags", []string{"error", "critical"})
	e := tagEntry("tags", "error, info")
	if f.Match(e) {
		t.Error("expected no match when a tag is missing")
	}
}

func TestTagFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := NewTagFilter("tags", []string{"ERROR"})
	e := tagEntry("tags", "error")
	if !f.Match(e) {
		t.Error("expected case-insensitive match")
	}
}

func TestTagFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewTagFilter("tags", []string{"error"})
	e := parser.LogEntry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Error("expected no match when field absent")
	}
}

func TestTagFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewTagFilter("tags", []string{"error"})
	e := parser.LogEntry{Fields: map[string]any{"tags": 42}}
	if f.Match(e) {
		t.Error("expected no match for non-string field")
	}
}
