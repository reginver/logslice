package cli

import (
	"testing"
)

func TestParseSuffixPairs_Valid(t *testing.T) {
	pairs, err := parseSuffixPairs([]string{"file=.go", "msg=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0] != [2]string{"file", ".go"} {
		t.Errorf("unexpected pair[0]: %v", pairs[0])
	}
	if pairs[1] != [2]string{"msg", "error"} {
		t.Errorf("unexpected pair[1]: %v", pairs[1])
	}
}

func TestParseSuffixPairs_Empty(t *testing.T) {
	pairs, err := parseSuffixPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected empty pairs")
	}
}

func TestParseSuffixPair_MissingParts(t *testing.T) {
	_, err := parseSuffixPair("noequals")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseSuffixPair_EmptyField(t *testing.T) {
	_, err := parseSuffixPair("=.go")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseSuffixPair_EmptySuffix(t *testing.T) {
	_, err := parseSuffixPair("file=")
	if err == nil {
		t.Fatal("expected error for empty suffix")
	}
}

func TestParseSuffixPair_ValueWithEquals(t *testing.T) {
	p, err := parseSuffixPair("msg=ends=here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != ([2]string{"msg", "ends=here"}) {
		t.Errorf("unexpected pair: %v", p)
	}
}
