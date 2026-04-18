package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]string, ts *time.Time) parser.LogEntry {
	return parser.LogEntry{Fields: fields, Timestamp: ts}
}

func TestCompositeFilter_AllMatch(t *testing.T) {
	f1, _ := filter.NewFieldFilter("level", "error")
	f2, _ := filter.NewFieldFilter("service", "auth")
	cf := filter.NewCompositeFilter(f1, f2)

	e := makeEntry(map[string]string{"level": "error", "service": "auth"}, nil)
	if !cf.Match(e) {
		t.Error("expected match when all filters pass")
	}
}

func TestCompositeFilter_OneFails(t *testing.T) {
	f1, _ := filter.NewFieldFilter("level", "error")
	f2, _ := filter.NewFieldFilter("service", "auth")
	cf := filter.NewCompositeFilter(f1, f2)

	e := makeEntry(map[string]string{"level": "error", "service": "payment"}, nil)
	if cf.Match(e) {
		t.Error("expected no match when one filter fails")
	}
}

func TestCompositeFilter_Empty(t *testing.T) {
	cf := filter.NewCompositeFilter()
	e := makeEntry(map[string]string{}, nil)
	if !cf.Match(e) {
		t.Error("empty composite filter should match everything")
	}
}

func TestAnyFilter_OneMatches(t *testing.T) {
	f1, _ := filter.NewFieldFilter("level", "error")
	f2, _ := filter.NewFieldFilter("level", "warn")
	af := filter.NewAnyFilter(f1, f2)

	e := makeEntry(map[string]string{"level": "warn"}, nil)
	if !af.Match(e) {
		t.Error("expected match when at least one filter passes")
	}
}

func TestAnyFilter_NoneMatch(t *testing.T) {
	f1, _ := filter.NewFieldFilter("level", "error")
	f2, _ := filter.NewFieldFilter("level", "warn")
	af := filter.NewAnyFilter(f1, f2)

	e := makeEntry(map[string]string{"level": "info"}, nil)
	if af.Match(e) {
		t.Error("expected no match when no filters pass")
	}
}

func TestAnyFilter_Empty(t *testing.T) {
	af := filter.NewAnyFilter()
	e := makeEntry(map[string]string{}, nil)
	if af.Match(e) {
		t.Error("empty any-filter should match nothing")
	}
}
