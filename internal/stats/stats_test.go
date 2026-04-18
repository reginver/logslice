package stats

import (
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	s := c.Summary()
	if s.TotalLines != 0 || s.MatchedLines != 0 || s.SkippedLines != 0 {
		t.Fatal("expected zero counters on new collector")
	}
	if s.FieldCounts == nil {
		t.Fatal("FieldCounts map should be initialised")
	}
}

func TestRecordTotal(t *testing.T) {
	c := NewCollector()
	c.RecordTotal()
	c.RecordTotal()
	if c.Summary().TotalLines != 2 {
		t.Fatalf("expected 2, got %d", c.Summary().TotalLines)
	}
}

func TestRecordMatched_UpdatesTimeBounds(t *testing.T) {
	c := NewCollector()
	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)

	c.RecordMatched(&t1)
	c.RecordMatched(&t2)
	c.RecordMatched(&t3)

	s := c.Summary()
	if s.MatchedLines != 3 {
		t.Fatalf("expected 3 matched, got %d", s.MatchedLines)
	}
	if !s.EarliestTime.Equal(t3) {
		t.Errorf("expected earliest %v, got %v", t3, s.EarliestTime)
	}
	if !s.LatestTime.Equal(t2) {
		t.Errorf("expected latest %v, got %v", t2, s.LatestTime)
	}
}

func TestRecordMatched_NilTimestamp(t *testing.T) {
	c := NewCollector()
	c.RecordMatched(nil)
	s := c.Summary()
	if s.MatchedLines != 1 {
		t.Fatalf("expected 1 matched, got %d", s.MatchedLines)
	}
	if s.EarliestTime != nil || s.LatestTime != nil {
		t.Error("time bounds should remain nil when timestamp is nil")
	}
}

func TestRecordField(t *testing.T) {
	c := NewCollector()
	c.RecordField("error")
	c.RecordField("info")
	c.RecordField("error")
	s := c.Summary()
	if s.FieldCounts["error"] != 2 {
		t.Errorf("expected error count 2, got %d", s.FieldCounts["error"])
	}
	if s.FieldCounts["info"] != 1 {
		t.Errorf("expected info count 1, got %d", s.FieldCounts["info"])
	}
}
