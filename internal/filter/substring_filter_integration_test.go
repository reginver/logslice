package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func TestSubstringFilter_WithComposite(t *testing.T) {
	sub, err := filter.NewSubstringFilter("msg", "timeout", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pfx, err := filter.NewPrefixFilter("svc", "api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	composite := filter.NewCompositeFilter(sub, pfx)

	match := parser.LogEntry{Fields: map[string]any{"msg": "connection timeout", "svc": "api-gateway"}}
	nomatch := parser.LogEntry{Fields: map[string]any{"msg": "connection timeout", "svc": "db-proxy"}}

	if !composite.Match(match) {
		t.Error("expected composite match")
	}
	if composite.Match(nomatch) {
		t.Error("expected composite no-match")
	}
}

func TestSubstringFilter_Negated(t *testing.T) {
	sub, err := filter.NewSubstringFilter("msg", "debug", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	not := filter.NewNotFilter(sub)

	noDebug := parser.LogEntry{Fields: map[string]any{"msg": "request received"}}
	withDebug := parser.LogEntry{Fields: map[string]any{"msg": "DEBUG: skipping"}}

	if !not.Match(noDebug) {
		t.Error("expected match for non-debug entry")
	}
	if not.Match(withDebug) {
		t.Error("expected no match for debug entry")
	}
}
