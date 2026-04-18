package pipeline_test

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/pipeline"
)

const sampleLogs = `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"startup"}
{"time":"2024-01-15T11:00:00Z","level":"error","msg":"crash"}
{"time":"2024-01-15T12:00:00Z","level":"info","msg":"restart"}
not-json
{"time":"2024-01-15T13:00:00Z","level":"debug","msg":"probe"}
`

func TestRun_NoFilters(t *testing.T) {
	res, err := pipeline.Run(strings.NewReader(sampleLogs), pipeline.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(res.Entries))
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestRun_TimeFilter(t *testing.T) {
	cfg := pipeline.Config{
		From: "2024-01-15T10:30:00Z",
		To:   "2024-01-15T12:30:00Z",
	}
	res, err := pipeline.Run(strings.NewReader(sampleLogs), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestRun_FieldFilter(t *testing.T) {
	cfg := pipeline.Config{
		Fields: map[string]string{"level": "error"},
	}
	res, err := pipeline.Run(strings.NewReader(sampleLogs), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(res.Entries))
	}
	if res.Entries[0].Fields["msg"] != "crash" {
		t.Errorf("unexpected entry: %+v", res.Entries[0])
	}
}

func TestRun_CombinedFilters(t *testing.T) {
	cfg := pipeline.Config{
		From:   "2024-01-15T09:00:00Z",
		To:     "2024-01-15T11:30:00Z",
		Fields: map[string]string{"level": "info"},
	}
	res, err := pipeline.Run(strings.NewReader(sampleLogs), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(res.Entries))
	}
}
