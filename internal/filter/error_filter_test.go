package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func errorEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewErrorFilter_Valid(t *testing.T) {
	f, err := filter.NewErrorFilter("error", []string{"timeout", "refused"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewErrorFilter_EmptyField(t *testing.T) {
	_, err := filter.NewErrorFilter("", []string{"timeout"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewErrorFilter_NoCodes(t *testing.T) {
	_, err := filter.NewErrorFilter("error", []string{})
	if err == nil {
		t.Fatal("expected error for empty codes")
	}
}

func TestNewErrorFilter_EmptyCode(t *testing.T) {
	_, err := filter.NewErrorFilter("error", []string{"timeout", ""})
	if err == nil {
		t.Fatal("expected error for empty code value")
	}
}

func TestErrorFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewErrorFilter("error", []string{"timeout"})
	e := errorEntry("error", "connection timeout reached")
	if !f.Match(e) {
		t.Fatal("expected match")
	}
}

func TestErrorFilter_Match_Miss(t *testing.T) {
	f, _ := filter.NewErrorFilter("error", []string{"refused"})
	e := errorEntry("error", "connection timeout reached")
	if f.Match(e) {
		t.Fatal("expected no match")
	}
}

func TestErrorFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := filter.NewErrorFilter("error", []string{"TIMEOUT"})
	e := errorEntry("error", "connection timeout reached")
	if !f.Match(e) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestErrorFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewErrorFilter("error", []string{"timeout"})
	e := parser.LogEntry{Fields: map[string]any{"msg": "hello"}}
	if f.Match(e) {
		t.Fatal("expected no match when field is absent")
	}
}

func TestErrorFilter_Match_MultipleCodesFirstHits(t *testing.T) {
	f, _ := filter.NewErrorFilter("error", []string{"refused", "timeout"})
	e := errorEntry("error", "dial timeout")
	if !f.Match(e) {
		t.Fatal("expected match on second code")
	}
}
