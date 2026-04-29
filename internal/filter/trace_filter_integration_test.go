package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func TestTraceFilter_WithComposite(t *testing.T) {
	trace, err := NewTraceFilter("trace_id", "abc*,xyz999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	level, err := NewLevelFilter("level", "warn", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	composite := NewCompositeFilter(trace, level)

	match := parser.Entry{Fields: map[string]any{
		"trace_id": "abc-0001",
		"level":    "error",
	}}
	noMatch := parser.Entry{Fields: map[string]any{
		"trace_id": "abc-0001",
		"level":    "debug",
	}}

	if !composite.Match(match) {
		t.Error("expected composite match")
	}
	if composite.Match(noMatch) {
		t.Error("expected composite no-match when level fails")
	}
}

func TestTraceFilter_Negated(t *testing.T) {
	trace, err := NewTraceFilter("trace_id", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	not := NewNotFilter(trace)

	matched := parser.Entry{Fields: map[string]any{"trace_id": "abc123"}}
	unmatched := parser.Entry{Fields: map[string]any{"trace_id": "xyz000"}}

	if not.Match(matched) {
		t.Error("expected negated filter to reject matched entry")
	}
	if !not.Match(unmatched) {
		t.Error("expected negated filter to pass unmatched entry")
	}
}
