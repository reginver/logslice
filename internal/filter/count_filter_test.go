package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func countEntry() parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{"msg": "hello"}}
}

func TestNewCountFilter_Valid(t *testing.T) {
	f, err := filter.NewCountFilter(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewCountFilter_ZeroInvalid(t *testing.T) {
	_, err := filter.NewCountFilter(0)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewCountFilter_NegativeInvalid(t *testing.T) {
	_, err := filter.NewCountFilter(-3)
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestCountFilter_PassesUpToMax(t *testing.T) {
	f, _ := filter.NewCountFilter(3)
	e := countEntry()
	for i := 0; i < 3; i++ {
		if !f.Match(e) {
			t.Fatalf("expected match on call %d", i+1)
		}
	}
}

func TestCountFilter_RejectsAfterMax(t *testing.T) {
	f, _ := filter.NewCountFilter(2)
	e := countEntry()
	f.Match(e)
	f.Match(e)
	if f.Match(e) {
		t.Fatal("expected no match after max reached")
	}
}

func TestCountFilter_ExactlyMax(t *testing.T) {
	f, _ := filter.NewCountFilter(1)
	e := countEntry()
	if !f.Match(e) {
		t.Fatal("expected first call to match")
	}
	if f.Match(e) {
		t.Fatal("expected second call to not match")
	}
}
