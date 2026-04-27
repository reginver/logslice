package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func TestKeywordFilter_WithComposite(t *testing.T) {
	kf, err := filter.NewKeywordFilter("msg", []string{"error", "fail"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lf, err := filter.NewFieldFilter("level", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cf := filter.NewCompositeFilter(kf, lf)

	hit := parser.LogEntry{Fields: map[string]any{"msg": "fatal error", "level": "error"}}
	if !cf.Match(hit) {
		t.Error("expected composite match")
	}

	miss := parser.LogEntry{Fields: map[string]any{"msg": "fatal error", "level": "info"}}
	if cf.Match(miss) {
		t.Error("expected composite miss when level does not match")
	}
}

func TestKeywordFilter_Negated(t *testing.T) {
	kf, err := filter.NewKeywordFilter("msg", []string{"debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nf, err := filter.NewNotFilter(kf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	noDebug := parser.LogEntry{Fields: map[string]any{"msg": "processing request"}}
	if !nf.Match(noDebug) {
		t.Error("expected negated match for non-debug entry")
	}

	debugEntry := parser.LogEntry{Fields: map[string]any{"msg": "debug: cache hit"}}
	if nf.Match(debugEntry) {
		t.Error("expected negated miss for debug entry")
	}
}
