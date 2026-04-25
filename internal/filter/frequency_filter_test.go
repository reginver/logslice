package filter

import (
	"testing"
)

func freqEntry(field, value string) map[string]any {
	return map[string]any{field: value}
}

func TestNewFrequencyFilter_Valid(t *testing.T) {
	f, err := NewFrequencyFilter("host", 2, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewFrequencyFilter_EmptyField(t *testing.T) {
	_, err := NewFrequencyFilter("", 1, -1)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewFrequencyFilter_ZeroMinCount(t *testing.T) {
	_, err := NewFrequencyFilter("host", 0, -1)
	if err == nil {
		t.Fatal("expected error for minCount=0")
	}
}

func TestNewFrequencyFilter_MaxLessThanMin(t *testing.T) {
	_, err := NewFrequencyFilter("host", 5, 2)
	if err == nil {
		t.Fatal("expected error when maxCount < minCount")
	}
}

func TestFrequencyFilter_Match_HitMinCount(t *testing.T) {
	f, _ := NewFrequencyFilter("host", 2, -1)
	entries := []map[string]any{
		freqEntry("host", "a"),
		freqEntry("host", "a"),
		freqEntry("host", "b"),
	}
	for _, e := range entries {
		f.Tally(e)
	}
	if !f.Match(freqEntry("host", "a")) {
		t.Error("expected 'a' to match (count=2 >= min=2)")
	}
	if f.Match(freqEntry("host", "b")) {
		t.Error("expected 'b' not to match (count=1 < min=2)")
	}
}

func TestFrequencyFilter_Match_MaxCount(t *testing.T) {
	f, _ := NewFrequencyFilter("host", 1, 2)
	for i := 0; i < 3; i++ {
		f.Tally(freqEntry("host", "x"))
	}
	if f.Match(freqEntry("host", "x")) {
		t.Error("expected 'x' not to match (count=3 > max=2)")
	}
}

func TestFrequencyFilter_Match_MissingField(t *testing.T) {
	f, _ := NewFrequencyFilter("host", 1, -1)
	if f.Match(map[string]any{"other": "val"}) {
		t.Error("expected no match for missing field")
	}
}

func TestFrequencyFilter_UnboundedMax(t *testing.T) {
	f, _ := NewFrequencyFilter("svc", 1, -1)
	for i := 0; i < 100; i++ {
		f.Tally(freqEntry("svc", "api"))
	}
	if !f.Match(freqEntry("svc", "api")) {
		t.Error("expected match with unbounded max")
	}
}
