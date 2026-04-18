package parser

import (
	"io"
	"strings"
	"testing"
)

func TestDetectFormat_JSON(t *testing.T) {
	input := strings.NewReader(`{"ts":"2024-01-01T00:00:00Z","msg":"hello"}` + "\n")
	fmt, r, err := DetectFormat(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != FormatJSON {
		t.Errorf("expected FormatJSON, got %v", fmt)
	}
	data, _ := io.ReadAll(r)
	if !strings.Contains(string(data), "hello") {
		t.Error("reader should still contain original content")
	}
}

func TestDetectFormat_Text(t *testing.T) {
	input := strings.NewReader("2024-01-01T00:00:00Z INFO booting\n")
	fmt, _, err := DetectFormat(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != FormatText {
		t.Errorf("expected FormatText, got %v", fmt)
	}
}

func TestDetectFormat_EmptyInput(t *testing.T) {
	input := strings.NewReader("")
	fmt, _, err := DetectFormat(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != FormatUnknown {
		t.Errorf("expected FormatUnknown, got %v", fmt)
	}
}

func TestDetectFormat_SkipsBlankLines(t *testing.T) {
	input := strings.NewReader("\n\n{\"msg\":\"hi\"}\n")
	fmt, _, err := DetectFormat(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != FormatJSON {
		t.Errorf("expected FormatJSON after blank lines, got %v", fmt)
	}
}
