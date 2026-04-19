package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func wildcardEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewWildcardFilter_Valid(t *testing.T) {
	_, err := filter.NewWildcardFilter("service", "auth-*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewWildcardFilter_EmptyField(t *testing.T) {
	_, err := filter.NewWildcardFilter("", "auth-*")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewWildcardFilter_EmptyPattern(t *testing.T) {
	_, err := filter.NewWildcardFilter("service", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestWildcardFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewWildcardFilter("service", "auth-*")
	e := wildcardEntry("service", "auth-service")
	if !f.Match(e) {
		t.Error("expected match")
	}
}

func TestWildcardFilter_Match_Miss(t *testing.T) {
	f, _ := filter.NewWildcardFilter("service", "auth-*")
	e := wildcardEntry("service", "billing-service")
	if f.Match(e) {
		t.Error("expected no match")
	}
}

func TestWildcardFilter_Match_MissingField(t *testing.T) {
	f, _ := filter.NewWildcardFilter("service", "auth-*")
	e := parser.LogEntry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Error("expected no match for missing field")
	}
}

func TestWildcardFilter_Match_QuestionMark(t *testing.T) {
	f, _ := filter.NewWildcardFilter("code", "50?")
	if !f.Match(wildcardEntry("code", "500")) {
		t.Error("expected match for 500")
	}
	if f.Match(wildcardEntry("code", "4000")) {
		t.Error("expected no match for 4000")
	}
}
