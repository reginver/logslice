package filter

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func alwaysMatch(_ parser.LogEntry) bool { return true }
func neverMatch(_ parser.LogEntry) bool  { return false }

type stubFilter struct{ result bool }

func (s *stubFilter) Match(_ parser.LogEntry) bool { return s.result }

func TestNewNotFilter_NilInner(t *testing.T) {
	_, err := NewNotFilter(nil)
	if err == nil {
		t.Fatal("expected error for nil inner filter")
	}
}

func TestNotFilter_NegatesTrue(t *testing.T) {
	f, err := NewNotFilter(&stubFilter{result: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Error("expected Match to return false when inner returns true")
	}
}

func TestNotFilter_NegatesFalse(t *testing.T) {
	f, err := NewNotFilter(&stubFilter{result: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if !f.Match(e) {
		t.Error("expected Match to return true when inner returns false")
	}
}

func TestNotFilter_WithFieldFilter(t *testing.T) {
	ff, err := NewFieldFilter("level", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nf, err := NewNotFilter(ff)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	errorEntry := parser.LogEntry{Fields: map[string]interface{}{"level": "error"}}
	infoEntry := parser.LogEntry{Fields: map[string]interface{}{"level": "info"}}

	if nf.Match(errorEntry) {
		t.Error("expected error entry to be excluded")
	}
	if !nf.Match(infoEntry) {
		t.Error("expected info entry to be included")
	}
}
