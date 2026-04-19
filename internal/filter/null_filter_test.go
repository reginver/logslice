package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func nullEntry(fields map[string]any) parser.LogEntry {
	return parser.LogEntry{Fields: fields}
}

func TestNewNullFilter_Valid(t *testing.T) {
	_, err := filter.NewNullFilter("error", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewNullFilter_EmptyField(t *testing.T) {
	_, err := filter.NewNullFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNullFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewNullFilter("error", false)
	e := nullEntry(map[string]any{"msg": "hello"})
	if !f.Match(e) {
		t.Error("expected match when field is absent")
	}
}

func TestNullFilter_Match_FieldNil(t *testing.T) {
	f, _ := filter.NewNullFilter("error", false)
	e := nullEntry(map[string]any{"error": nil})
	if !f.Match(e) {
		t.Error("expected match when field is nil")
	}
}

func TestNullFilter_Match_FieldEmpty(t *testing.T) {
	f, _ := filter.NewNullFilter("error", false)
	e := nullEntry(map[string]any{"error": ""})
	if !f.Match(e) {
		t.Error("expected match when field is empty string")
	}
}

func TestNullFilter_Match_FieldPresent(t *testing.T) {
	f, _ := filter.NewNullFilter("error", false)
	e := nullEntry(map[string]any{"error": "something went wrong"})
	if f.Match(e) {
		t.Error("expected no match when field has value")
	}
}

func TestNullFilter_Inverted_MatchesNonNull(t *testing.T) {
	f, _ := filter.NewNullFilter("error", true)
	e := nullEntry(map[string]any{"error": "oops"})
	if !f.Match(e) {
		t.Error("expected match when field is present and inverted")
	}
}

func TestNullFilter_Inverted_NoMatchOnAbsent(t *testing.T) {
	f, _ := filter.NewNullFilter("error", true)
	e := nullEntry(map[string]any{"msg": "ok"})
	if f.Match(e) {
		t.Error("expected no match when field absent and inverted")
	}
}
