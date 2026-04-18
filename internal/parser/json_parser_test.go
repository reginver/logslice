package parser

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestParseJSON_ValidEntries(t *testing.T) {
	input := `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"started"}
{"time":"2024-01-15T10:01:00Z","level":"error","msg":"failed"}`

	entries, err := ParseJSON(strings.NewReader(input), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entries[0].Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entries[0].Timestamp)
	}
	if entries[1].Fields["level"] != "error" {
		t.Errorf("expected level=error, got %v", entries[1].Fields["level"])
	}
}

func TestParseJSON_SkipsInvalidLines(t *testing.T) {
	input := `{"time":"2024-01-15T10:00:00Z","msg":"ok"}
not json at all
{"time":"2024-01-15T10:02:00Z","msg":"also ok"}`

	var errBuf bytes.Buffer
	entries, err := ParseJSON(strings.NewReader(input), &errBuf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !strings.Contains(errBuf.String(), "skipping invalid JSON") {
		t.Errorf("expected warning in errOut, got: %s", errBuf.String())
	}
}

func TestParseJSON_UnixTimestamp(t *testing.T) {
	input := `{"time":1705312800,"msg":"unix ts"}`

	entries, err := ParseJSON(strings.NewReader(input), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp for unix epoch value")
	}
}

func TestParseJSON_EmptyInput(t *testing.T) {
	entries, err := ParseJSON(strings.NewReader(""), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
