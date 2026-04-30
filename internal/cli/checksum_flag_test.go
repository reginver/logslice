package cli

import (
	"testing"
)

func TestParseChecksumPairs_Valid(t *testing.T) {
	filters, err := parseChecksumPairs([]string{"id=md5:deadbeef", "body=sha256:cafe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseChecksumPairs_Empty(t *testing.T) {
	filters, err := parseChecksumPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseChecksumPair_MissingEquals(t *testing.T) {
	_, err := parseChecksumPair("idmd5deadbeef")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseChecksumPair_EmptyField(t *testing.T) {
	_, err := parseChecksumPair("=md5:deadbeef")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseChecksumPair_MissingColon(t *testing.T) {
	_, err := parseChecksumPair("id=md5deadbeef")
	if err == nil {
		t.Fatal("expected error for missing ':'")
	}
}

func TestParseChecksumPair_EmptyAlgo(t *testing.T) {
	_, err := parseChecksumPair("id=:deadbeef")
	if err == nil {
		t.Fatal("expected error for empty algo")
	}
}

func TestParseChecksumPair_EmptyPrefix(t *testing.T) {
	_, err := parseChecksumPair("id=md5:")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestParseChecksumPair_UnknownAlgo(t *testing.T) {
	_, err := parseChecksumPair("id=crc32:abcd")
	if err == nil {
		t.Fatal("expected error for unknown algo")
	}
}
