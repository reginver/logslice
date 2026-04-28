package cli

import (
	"testing"
)

func TestParseTagPairs_Valid(t *testing.T) {
	filters, err := parseTagPairs([]string{"tags=error,critical", "labels=prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseTagPairs_Empty(t *testing.T) {
	filters, err := parseTagPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseTagPair_MissingEquals(t *testing.T) {
	_, err := parseTagPair("tagserror")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseTagPair_EmptyField(t *testing.T) {
	_, err := parseTagPair("=error")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseTagPair_EmptyTagList(t *testing.T) {
	_, err := parseTagPair("tags=")
	if err == nil {
		t.Fatal("expected error for empty tag list")
	}
}

func TestParseTagPair_SingleTag(t *testing.T) {
	f, err := parseTagPair("tags=error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
