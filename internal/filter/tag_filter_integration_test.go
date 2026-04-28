package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func TestTagFilter_WithComposite(t *testing.T) {
	tag, err := NewTagFilter("tags", []string{"error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exists, err := NewExistsFilter("request_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	composite := NewCompositeFilter(tag, exists)

	hit := parser.LogEntry{Fields: map[string]any{
		"tags":       "error,warning",
		"request_id": "abc-123",
	}}
	miss := parser.LogEntry{Fields: map[string]any{
		"tags": "error,warning",
	}}

	if !composite.Match(hit) {
		t.Error("expected composite match on hit entry")
	}
	if composite.Match(miss) {
		t.Error("expected composite no-match on miss entry")
	}
}

func TestTagFilter_Negated(t *testing.T) {
	tag, err := NewTagFilter("tags", []string{"debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	not, err := NewNotFilter(tag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hasDebug := parser.LogEntry{Fields: map[string]any{"tags": "debug,info"}}
	noDebug := parser.LogEntry{Fields: map[string]any{"tags": "info,warning"}}

	if not.Match(hasDebug) {
		t.Error("expected negated filter to reject entry with debug tag")
	}
	if !not.Match(noDebug) {
		t.Error("expected negated filter to pass entry without debug tag")
	}
}
