package cli

import (
	"testing"
)

func TestParsePrefixPairs_Valid(t *testing.T) {
	pairs, err := parsePrefixPairs([]string{"msg:ERROR", "service:auth"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].Field != "msg" || pairs[0].Prefix != "ERROR" {
		t.Errorf("unexpected pair[0]: %+v", pairs[0])
	}
	if pairs[1].Field != "service" || pairs[1].Prefix != "auth" {
		t.Errorf("unexpected pair[1]: %+v", pairs[1])
	}
}

func TestParsePrefixPairs_Empty(t *testing.T) {
	pairs, err := parsePrefixPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected empty slice, got %d", len(pairs))
	}
}

func TestParsePrefixPair_MissingParts(t *testing.T) {
	_, err := parsePrefixPair("nodcolon")
	if err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestParsePrefixPair_EmptyField(t *testing.T) {
	_, err := parsePrefixPair(":value")
	if err == nil {
		t.Error("expected error for empty field")
	}
}

func TestParsePrefixPair_EmptyPrefix(t *testing.T) {
	_, err := parsePrefixPair("field:")
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}

func TestParsePrefixPair_Valid(t *testing.T) {
	p, err := parsePrefixPair("level:WA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Field != "level" || p.Prefix != "WA" {
		t.Errorf("unexpected result: %+v", p)
	}
}
