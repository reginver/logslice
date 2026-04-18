package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParseText_ValidEntries(t *testing.T) {
	input := strings.NewReader(
		"2024-03-01T10:00:00Z INFO user logged in\n" +
			"2024-03-01T10:01:00Z ERROR disk full\n",
	)
	entries, err := ParseText(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Level != "INFO" {
		t.Errorf("expected INFO, got %s", entries[0].Level)
	}
	if entries[1].Message != "disk full" {
		t.Errorf("unexpected message: %s", entries[1].Message)
	}
}

func TestParseText_TimestampParsed(t *testing.T) {
	input := strings.NewReader("2024-06-15T08:30:00Z WARN high memory\n")
	entries, err := ParseText(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	expected, _ := time.Parse(time.RFC3339, "2024-06-15T08:30:00Z")
	if entries[0].Timestamp == nil || !entries[0].Timestamp.Equal(expected) {
		t.Errorf("timestamp mismatch: got %v", entries[0].Timestamp)
	}
}

func TestParseText_UnmatchedLineRetainsRaw(t *testing.T) {
	input := strings.NewReader("this is not a structured log line\n")
	entries, err := ParseText(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if entries[0].Timestamp != nil {
		t.Error("expected nil timestamp for unmatched line")
	}
	if entries[0].Raw != "this is not a structured log line" {
		t.Errorf("unexpected raw: %s", entries[0].Raw)
	}
}

func TestParseText_EmptyInput(t *testing.T) {
	entries, err := ParseText(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseText_SkipsBlankLines(t *testing.T) {
	input := strings.NewReader("\n2024-01-01T00:00:00Z DEBUG boot\n\n")
	entries, err := ParseText(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}
