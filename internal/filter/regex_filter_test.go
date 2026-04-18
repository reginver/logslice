package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func entry(fields map[string]any) parser.LogEntry {
	return parser.LogEntry{Fields: fields}
}

func TestNewRegexFilter_Valid(t *testing.T) {
	f, err := NewRegexFilter("level", `^(error|warn)$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewRegexFilter_EmptyField(t *testing.T) {
	_, err := NewRegexFilter("", `.*`)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRegexFilter_InvalidPattern(t *testing.T) {
	_, err := NewRegexFilter("msg", `[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRegexFilter_Match_Hit(t *testing.T) {
	f, _ := NewRegexFilter("level", `^error$`)
	if !f.Match(entry(map[string]any{"level": "error"})) {
		t.Error("expected match")
	}
}

func TestRegexFilter_Match_Miss(t *testing.T) {
	f, _ := NewRegexFilter("level", `^error$`)
	if f.Match(entry(map[string]any{"level": "info"})) {
		t.Error("expected no match")
	}
}

func TestRegexFilter_Match_MissingField(t *testing.T) {
	f, _ := NewRegexFilter("level", `.*`)
	if f.Match(entry(map[string]any{"msg": "hello"})) {
		t.Error("expected no match when field absent")
	}
}

func TestRegexFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewRegexFilter("code", `^404$`)
	if !f.Match(entry(map[string]any{"code": 404})) {
		t.Error("expected match on numeric field formatted as string")
	}
}
