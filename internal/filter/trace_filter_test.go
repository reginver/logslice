package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func traceEntry(field, value string) parser.Entry {
	return parser.Entry{Fields: map[string]any{field: value}}
}

func TestNewTraceFilter_Valid(t *testing.T) {
	f, err := NewTraceFilter("trace_id", "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewTraceFilter_EmptyField(t *testing.T) {
	_, err := NewTraceFilter("", "abc")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTraceFilter_EmptySpec(t *testing.T) {
	_, err := NewTraceFilter("trace_id", "")
	if err == nil {
		t.Fatal("expected error for empty spec")
	}
}

func TestNewTraceFilter_EmptyIDInSpec(t *testing.T) {
	_, err := NewTraceFilter("trace_id", "abc,,def")
	if err == nil {
		t.Fatal("expected error for empty ID in spec")
	}
}

func TestTraceFilter_Match_ExactHit(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc123")
	if !f.Match(traceEntry("trace_id", "abc123")) {
		t.Error("expected match on exact ID")
	}
}

func TestTraceFilter_Match_ExactMiss(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc123")
	if f.Match(traceEntry("trace_id", "xyz999")) {
		t.Error("expected no match")
	}
}

func TestTraceFilter_Match_PrefixHit(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc*")
	if !f.Match(traceEntry("trace_id", "abc-0001")) {
		t.Error("expected prefix match")
	}
}

func TestTraceFilter_Match_PrefixMiss(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc*")
	if f.Match(traceEntry("trace_id", "xyz-0001")) {
		t.Error("expected no match")
	}
}

func TestTraceFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc123")
	e := parser.Entry{Fields: map[string]any{"other": "abc123"}}
	if f.Match(e) {
		t.Error("expected no match when field absent")
	}
}

func TestTraceFilter_Match_NonStringValue(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc123")
	e := parser.Entry{Fields: map[string]any{"trace_id": 42}}
	if f.Match(e) {
		t.Error("expected no match for non-string value")
	}
}

func TestTraceFilter_Match_MultipleSpecs(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", "abc123,def*,xyz789")
	if !f.Match(traceEntry("trace_id", "abc123")) {
		t.Error("expected exact match")
	}
	if !f.Match(traceEntry("trace_id", "def-0042")) {
		t.Error("expected prefix match")
	}
	if !f.Match(traceEntry("trace_id", "xyz789")) {
		t.Error("expected exact match for xyz789")
	}
	if f.Match(traceEntry("trace_id", "nope")) {
		t.Error("expected no match for unrelated ID")
	}
}
