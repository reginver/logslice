package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func existsEntry(fields map[string]any) parser.LogEntry {
	return parser.LogEntry{Fields: fields}
}

func TestNewExistsFilter_Valid(t *testing.T) {
	_, err := filter.NewExistsFilter("level", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewExistsFilter_EmptyField(t *testing.T) {
	_, err := filter.NewExistsFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestExistsFilter_Match_FieldPresent(t *testing.T) {
	f, _ := filter.NewExistsFilter("level", false)
	e := existsEntry(map[string]any{"level": "info"})
	if !f.Match(e) {
		t.Error("expected match when field is present")
	}
}

func TestExistsFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewExistsFilter("level", false)
	e := existsEntry(map[string]any{"msg": "hello"})
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestExistsFilter_Negate_FieldAbsent(t *testing.T) {
	f, _ := filter.NewExistsFilter("level", true)
	e := existsEntry(map[string]any{"msg": "hello"})
	if !f.Match(e) {
		t.Error("expected match when field absent and negate=true")
	}
}

func TestExistsFilter_Negate_FieldPresent(t *testing.T) {
	f, _ := filter.NewExistsFilter("level", true)
	e := existsEntry(map[string]any{"level": "error"})
	if f.Match(e) {
		t.Error("expected no match when field present and negate=true")
	}
}

func TestExistsFilter_NilFields(t *testing.T) {
	f, _ := filter.NewExistsFilter("level", false)
	e := parser.LogEntry{}
	if f.Match(e) {
		t.Error("expected no match for entry with nil fields")
	}
}
