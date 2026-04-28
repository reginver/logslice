package cli

import (
	"testing"
)

func TestParseErrorPairs_Valid(t *testing.T) {
	filters, err := parseErrorPairs([]string{"error=timeout,refused", "cause=oom"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseErrorPairs_Empty(t *testing.T) {
	filters, err := parseErrorPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseErrorPair_MissingEquals(t *testing.T) {
	_, err := parseErrorPair("errortimeout")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseErrorPair_EmptyField(t *testing.T) {
	_, err := parseErrorPair("=timeout")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseErrorPair_EmptyCodeList(t *testing.T) {
	_, err := parseErrorPair("error=")
	if err == nil {
		t.Fatal("expected error for empty code list")
	}
}

func TestParseErrorPair_SingleCode(t *testing.T) {
	f, err := parseErrorPair("error=timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
