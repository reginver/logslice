package cli

import (
	"testing"
)

func TestParseRangePairs_Valid(t *testing.T) {
	pairs, err := parseRangePairs([]string{"latency:10:500", "code:200:299"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].Field != "latency" || pairs[0].Min != 10 || pairs[0].Max != 500 {
		t.Errorf("unexpected first pair: %+v", pairs[0])
	}
	if pairs[1].Field != "code" || pairs[1].Min != 200 || pairs[1].Max != 299 {
		t.Errorf("unexpected second pair: %+v", pairs[1])
	}
}

func TestParseRangePairs_Empty(t *testing.T) {
	pairs, err := parseRangePairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestParseRangePair_MissingParts(t *testing.T) {
	_, err := parseRangePair("latency:10")
	if err == nil {
		t.Fatal("expected error for missing max")
	}
}

func TestParseRangePair_EmptyField(t *testing.T) {
	_, err := parseRangePair(":10:100")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseRangePair_InvalidMin(t *testing.T) {
	_, err := parseRangePair("score:abc:100")
	if err == nil {
		t.Fatal("expected error for non-numeric min")
	}
}

func TestParseRangePair_InvalidMax(t *testing.T) {
	_, err := parseRangePair("score:0:xyz")
	if err == nil {
		t.Fatal("expected error for non-numeric max")
	}
}

func TestParseRangePair_FloatValues(t *testing.T) {
	p, err := parseRangePair("ratio:0.1:0.9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Min != 0.1 || p.Max != 0.9 {
		t.Errorf("unexpected values: %+v", p)
	}
}
