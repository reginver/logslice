package cli

import (
	"testing"
)

func TestParseContainsPairs_Valid(t *testing.T) {
	pairs, err := parseContainsPairs([]string{"level=error,warn", "env=prod,staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].Field != "level" {
		t.Errorf("expected field 'level', got %q", pairs[0].Field)
	}
	if len(pairs[0].Values) != 2 || pairs[0].Values[0] != "error" || pairs[0].Values[1] != "warn" {
		t.Errorf("unexpected values: %v", pairs[0].Values)
	}
}

func TestParseContainsPairs_Empty(t *testing.T) {
	pairs, err := parseContainsPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected empty, got %d", len(pairs))
	}
}

func TestParseContainsPair_MissingParts(t *testing.T) {
	_, err := parseContainsPair("levelonly")
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseContainsPair_EmptyField(t *testing.T) {
	_, err := parseContainsPair("=error,warn")
	if err == nil {
		t.Error("expected error for empty field")
	}
}

func TestParseContainsPair_EmptyValues(t *testing.T) {
	_, err := parseContainsPair("level=")
	if err == nil {
		t.Error("expected error for empty values")
	}
}

func TestParseContainsPair_SingleValue(t *testing.T) {
	p, err := parseContainsPair("status=200")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Field != "status" {
		t.Errorf("expected field 'status', got %q", p.Field)
	}
	if len(p.Values) != 1 || p.Values[0] != "200" {
		t.Errorf("unexpected values: %v", p.Values)
	}
}
