package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func TestPathFilter_WithComposite(t *testing.T) {
	pf, err := filter.NewPathFilter("url", "/api/*")
	if err != nil {
		t.Fatalf("NewPathFilter: %v", err)
	}
	ff, err := filter.NewFieldFilter("method", "GET")
	if err != nil {
		t.Fatalf("NewFieldFilter: %v", err)
	}
	comp := filter.NewCompositeFilter(pf, ff)

	hit := parser.LogEntry{Fields: map[string]any{"url": "/api/users", "method": "GET"}}
	miss := parser.LogEntry{Fields: map[string]any{"url": "/api/users", "method": "POST"}}

	if !comp.Match(hit) {
		t.Error("expected composite match")
	}
	if comp.Match(miss) {
		t.Error("expected composite miss")
	}
}

func TestPathFilter_Negated(t *testing.T) {
	pf, err := filter.NewPathFilter("url", "/internal/*")
	if err != nil {
		t.Fatalf("NewPathFilter: %v", err)
	}
	not := filter.NewNotFilter(pf)

	public := parser.LogEntry{Fields: map[string]any{"url": "/api/data"}}
	internal := parser.LogEntry{Fields: map[string]any{"url": "/internal/metrics"}}

	if !not.Match(public) {
		t.Error("expected negated match for public path")
	}
	if not.Match(internal) {
		t.Error("expected negated miss for internal path")
	}
}
