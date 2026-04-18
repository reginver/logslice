package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_BasicPipeline(t *testing.T) {
	input := writeTempLog(t, `{"timestamp":"2024-01-01T10:00:00Z","level":"info","msg":"hello"}`+"\n")
	output := filepath.Join(t.TempDir(), "out.log")

	err := Run([]string{"-input", input, "-output", output, "-format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "hello") {
		t.Errorf("expected output to contain 'hello', got: %s", data)
	}
}

func TestRun_InvalidTimeFilter(t *testing.T) {
	input := writeTempLog(t, `{"timestamp":"2024-01-01T10:00:00Z","level":"info"}`+"\n")
	err := Run([]string{"-input", input, "-from", "not-a-time"})
	if err == nil {
		t.Error("expected error for invalid time filter")
	}
}

func TestRun_MissingInputFile(t *testing.T) {
	err := Run([]string{"-input", "/nonexistent/path/file.log"})
	if err == nil {
		t.Error("expected error for missing input file")
	}
}
