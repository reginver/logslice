package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func sampleEntry(msg string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{"msg": msg}}
}

func TestNewSampleFilter_Valid(t *testing.T) {
	f, err := NewSampleFilter(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewSampleFilter_ZeroInvalid(t *testing.T) {
	_, err := NewSampleFilter(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestSampleFilter_N1_PassesAll(t *testing.T) {
	f, _ := NewSampleFilter(1)
	for i := 0; i < 10; i++ {
		if !f.Match(sampleEntry("x")) {
			t.Errorf("expected entry %d to pass with n=1", i)
		}
	}
}

func TestSampleFilter_N3_PassesEveryThird(t *testing.T) {
	f, _ := NewSampleFilter(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = f.Match(sampleEntry("x"))
	}
	// Positions 2, 5, 8 (0-indexed) should be true
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("entry %d: got %v, want %v", i, got, expected[i])
		}
	}
}

func TestSampleFilter_N2_PassesEverySecond(t *testing.T) {
	f, _ := NewSampleFilter(2)
	passed := 0
	for i := 0; i < 10; i++ {
		if f.Match(sampleEntry("x")) {
			passed++
		}
	}
	if passed != 5 {
		t.Errorf("expected 5 passed, got %d", passed)
	}
}
