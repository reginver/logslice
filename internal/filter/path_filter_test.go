package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func pathEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewPathFilter_Valid(t *testing.T) {
	_, err := filter.NewPathFilter("url", "/api/*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewPathFilter_EmptyField(t *testing.T) {
	_, err := filter.NewPathFilter("", "/api/*")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewPathFilter_EmptyPattern(t *testing.T) {
	_, err := filter.NewPathFilter("url", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewPathFilter_InvalidPattern(t *testing.T) {
	_, err := filter.NewPathFilter("url", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid glob pattern")
	}
}

func TestPathFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewPathFilter("url", "/api/v1/*")
	e := pathEntry("url", "/api/v1/users")
	if !f.Match(e) {
		t.Error("expected match")
	}
}

func TestPathFilter_Match_Miss(t *testing.T) {
	f, _ := filter.NewPathFilter("url", "/api/v1/*")
	e := pathEntry("url", "/api/v2/users")
	if f.Match(e) {
		t.Error("expected no match")
	}
}

func TestPathFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewPathFilter("url", "/api/*")
	e := parser.LogEntry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Error("expected no match for absent field")
	}
}

func TestPathFilter_Match_ExactPath(t *testing.T) {
	f, _ := filter.NewPathFilter("url", "/health")
	e := pathEntry("url", "/health")
	if !f.Match(e) {
		t.Error("expected match for exact path")
	}
}

func TestPathFilter_Match_Wildcard_Root(t *testing.T) {
	f, _ := filter.NewPathFilter("url", "/*")
	e := pathEntry("url", "/ping")
	if !f.Match(e) {
		t.Error("expected match")
	}
}
