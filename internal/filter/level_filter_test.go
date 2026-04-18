package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func levelEntry(level string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{"level": level}}
}

func TestNewLevelFilter_Valid(t *testing.T) {
	_, err := NewLevelFilter("level", "info", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLevelFilter_EmptyField(t *testing.T) {
	_, err := NewLevelFilter("", "info", "error")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewLevelFilter_UnknownMin(t *testing.T) {
	_, err := NewLevelFilter("level", "verbose", "error")
	if err == nil {
		t.Fatal("expected error for unknown min level")
	}
}

func TestNewLevelFilter_MinAboveMax(t *testing.T) {
	_, err := NewLevelFilter("level", "error", "info")
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestLevelFilter_Match_InRange(t *testing.T) {
	f, _ := NewLevelFilter("level", "info", "error")
	for _, lvl := range []string{"info", "warn", "error", "INFO", "WARN"} {
		if !f.Match(levelEntry(lvl)) {
			t.Errorf("expected match for level %q", lvl)
		}
	}
}

func TestLevelFilter_Match_OutOfRange(t *testing.T) {
	f, _ := NewLevelFilter("level", "info", "error")
	for _, lvl := range []string{"debug", "trace", "fatal"} {
		if f.Match(levelEntry(lvl)) {
			t.Errorf("expected no match for level %q", lvl)
		}
	}
}

func TestLevelFilter_Match_MissingField(t *testing.T) {
	f, _ := NewLevelFilter("level", "info", "error")
	e := parser.LogEntry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Fatal("expected no match when field absent")
	}
}

func TestLevelFilter_Match_OpenBounds(t *testing.T) {
	f, _ := NewLevelFilter("level", "", "")
	for _, lvl := range []string{"trace", "debug", "info", "warn", "error", "fatal"} {
		if !f.Match(levelEntry(lvl)) {
			t.Errorf("expected match for level %q with open bounds", lvl)
		}
	}
}
