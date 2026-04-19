package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func TestDateFilter_WithComposite(t *testing.T) {
	ts1 := time.Date(2024, 5, 10, 9, 0, 0, 0, time.UTC)
	ts2 := time.Date(2024, 5, 11, 9, 0, 0, 0, time.UTC)

	df, _ := NewDateFilter("2024-05-10", "")
	lf, _ := NewLevelFilter("level", "warn", "error")

	cf := NewCompositeFilter(df, lf)

	match := parser.Entry{Timestamp: &ts1, Fields: map[string]any{"level": "error"}}
	nomatch1 := parser.Entry{Timestamp: &ts2, Fields: map[string]any{"level": "error"}}
	nomatch2 := parser.Entry{Timestamp: &ts1, Fields: map[string]any{"level": "debug"}}

	if !cf.Match(match) {
		t.Fatal("expected composite match")
	}
	if cf.Match(nomatch1) {
		t.Fatal("expected no match: wrong date")
	}
	if cf.Match(nomatch2) {
		t.Fatal("expected no match: wrong level")
	}
}

func TestDateFilter_Negated(t *testing.T) {
	ts := time.Date(2024, 5, 10, 9, 0, 0, 0, time.UTC)
	other := time.Date(2024, 5, 11, 9, 0, 0, 0, time.UTC)

	df, _ := NewDateFilter("2024-05-10", "")
	nf, _ := NewNotFilter(df)

	if nf.Match(parser.Entry{Timestamp: &ts}) {
		t.Fatal("expected no match on negated matching date")
	}
	if !nf.Match(parser.Entry{Timestamp: &other}) {
		t.Fatal("expected match on negated non-matching date")
	}
}
