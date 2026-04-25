package filter

import (
	"testing"
)

func TestFrequencyFilter_WithComposite(t *testing.T) {
	freq, err := NewFrequencyFilter("env", 2, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries := []map[string]any{
		{"env": "prod", "level": "error"},
		{"env": "prod", "level": "warn"},
		{"env": "dev", "level": "info"},
	}
	for _, e := range entries {
		freq.Tally(e)
	}

	level, err := NewLevelFilter("level", "warn", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	composite := NewCompositeFilter(freq, level)

	// prod+error: freq(prod)=2 ✓, level=error ✓
	if !composite.Match(map[string]any{"env": "prod", "level": "error"}) {
		t.Error("expected prod+error to match composite")
	}
	// dev+info: freq(dev)=1 < 2 ✗
	if composite.Match(map[string]any{"env": "dev", "level": "info"}) {
		t.Error("expected dev+info not to match composite")
	}
	// prod+info: level=info below warn ✗
	if composite.Match(map[string]any{"env": "prod", "level": "info"}) {
		t.Error("expected prod+info not to match composite")
	}
}

func TestFrequencyFilter_Negated(t *testing.T) {
	freq, err := NewFrequencyFilter("user", 3, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 5; i++ {
		freq.Tally(map[string]any{"user": "alice"})
	}
	freq.Tally(map[string]any{"user": "bob"})

	not := NewNotFilter(freq)

	// alice appears 5 times — freq matches, NOT inverts → false
	if not.Match(map[string]any{"user": "alice"}) {
		t.Error("expected alice not to match negated filter")
	}
	// bob appears 1 time — freq does not match, NOT inverts → true
	if !not.Match(map[string]any{"user": "bob"}) {
		t.Error("expected bob to match negated filter")
	}
}
