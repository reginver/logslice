package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func TestWildcardFilter_WithComposite(t *testing.T) {
	wf, err := filter.NewWildcardFilter("service", "auth-*")
	if err != nil {
		t.Fatal(err)
	}
	lf, err := filter.NewLevelFilter("level", "warn", "error")
	if err != nil {
		t.Fatal(err)
	}
	cf := filter.NewCompositeFilter(wf, lf)

	hit := parser.LogEntry{Fields: map[string]any{"service": "auth-api", "level": "error"}}
	if !cf.Match(hit) {
		t.Error("expected composite match")
	}

	miss := parser.LogEntry{Fields: map[string]any{"service": "auth-api", "level": "debug"}}
	if cf.Match(miss) {
		t.Error("expected composite miss on low level")
	}
}

func TestWildcardFilter_Negated(t *testing.T) {
	wf, err := filter.NewWildcardFilter("env", "prod-*")
	if err != nil {
		t.Fatal(err)
	}
	nf, err := filter.NewNotFilter(wf)
	if err != nil {
		t.Fatal(err)
	}

	prodEntry := parser.LogEntry{Fields: map[string]any{"env": "prod-us"}}
	if nf.Match(prodEntry) {
		t.Error("expected negated filter to reject prod entry")
	}

	devEntry := parser.LogEntry{Fields: map[string]any{"env": "dev-local"}}
	if !nf.Match(devEntry) {
		t.Error("expected negated filter to pass dev entry")
	}
}
