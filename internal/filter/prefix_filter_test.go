package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func prefixEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{
		Fields: map[string]interface{}{field: value},
	}
}

func TestNewPrefixFilter_Valid(t *testing.T) {
	f, err := NewPrefixFilter("service", "auth")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewPrefixFilter_EmptyField(t *testing.T) {
	_, err := NewPrefixFilter("", "auth")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewPrefixFilter_EmptyPrefix(t *testing.T) {
	_, err := NewPrefixFilter("service", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestPrefixFilter_Match_Hit(t *testing.T) {
	f, _ := NewPrefixFilter("service", "auth")
	e := prefixEntry("service", "auth-service")
	if !f.Match(e) {
		t.Error("expected match")
	}
}

func TestPrefixFilter_Match_Miss(t *testing.T) {
	f, _ := NewPrefixFilter("service", "auth")
	e := prefixEntry("service", "billing-service")
	if f.Match(e) {
		t.Error("expected no match")
	}
}

func TestPrefixFilter_Match_MissingField(t *testing.T) {
	f, _ := NewPrefixFilter("service", "auth")
	e := parser.LogEntry{Fields: map[string]interface{}{"other": "auth-x"}}
	if f.Match(e) {
		t.Error("expected no match for missing field")
	}
}

func TestPrefixFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewPrefixFilter("code", "40")
	e := parser.LogEntry{Fields: map[string]interface{}{"code": 404}}
	if f.Match(e) {
		t.Error("expected no match for non-string field")
	}
}

func TestPrefixFilter_Match_ExactPrefix(t *testing.T) {
	f, _ := NewPrefixFilter("service", "auth")
	e := prefixEntry("service", "auth")
	if !f.Match(e) {
		t.Error("expected match when value equals prefix")
	}
}
