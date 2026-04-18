package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_BasicCounters(t *testing.T) {
	s := Summary{
		TotalLines:   10,
		MatchedLines: 7,
		SkippedLines: 3,
		FieldCounts:  map[string]int{},
	}
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	for _, want := range []string{"Total lines:", "10", "Matched lines:", "7", "Skipped lines:", "3"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nGot:\n%s", want, out)
		}
	}
}

func TestPrint_TimeBounds(t *testing.T) {
	t1 := time.Date(2024, 3, 15, 9, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 3, 15, 17, 30, 0, 0, time.UTC)
	s := Summary{
		FieldCounts:  map[string]int{},
		EarliestTime: &t1,
		LatestTime:   &t2,
	}
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "2024-03-15T09:00:00Z") {
		t.Errorf("expected earliest time in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2024-03-15T17:30:00Z") {
		t.Errorf("expected latest time in output, got:\n%s", out)
	}
}

func TestPrint_FieldCounts(t *testing.T) {
	s := Summary{
		FieldCounts: map[string]int{"error": 5, "warn": 2},
	}
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "error") || !strings.Contains(out, "5") {
		t.Errorf("expected error field count in output, got:\n%s", out)
	}
	if !strings.Contains(out, "warn") || !strings.Contains(out, "2") {
		t.Errorf("expected warn field count in output, got:\n%s", out)
	}
}

func TestPrint_NoTimeBounds(t *testing.T) {
	s := Summary{FieldCounts: map[string]int{}}
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()
	if strings.Contains(out, "Earliest") || strings.Contains(out, "Latest") {
		t.Errorf("should not print time bounds when nil, got:\n%s", out)
	}
}
