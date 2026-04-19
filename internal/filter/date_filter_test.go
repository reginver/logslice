package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func dateEntry(ts *time.Time, fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: ts, Fields: fields}
}

func TestNewDateFilter_Valid(t *testing.T) {
	_, err := NewDateFilter("2024-03-15", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDateFilter_EmptyDate(t *testing.T) {
	_, err := NewDateFilter("", "")
	if err == nil {
		t.Fatal("expected error for empty date")
	}
}

func TestNewDateFilter_InvalidFormat(t *testing.T) {
	_, err := NewDateFilter("15-03-2024", "")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestDateFilter_Match_Hit(t *testing.T) {
	ts := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	f, _ := NewDateFilter("2024-03-15", "")
	if !f.Match(dateEntry(&ts, nil)) {
		t.Fatal("expected match")
	}
}

func TestDateFilter_Match_Miss(t *testing.T) {
	ts := time.Date(2024, 3, 16, 10, 30, 0, 0, time.UTC)
	f, _ := NewDateFilter("2024-03-15", "")
	if f.Match(dateEntry(&ts, nil)) {
		t.Fatal("expected no match")
	}
}

func TestDateFilter_Match_NilTimestamp(t *testing.T) {
	f, _ := NewDateFilter("2024-03-15", "")
	if f.Match(dateEntry(nil, nil)) {
		t.Fatal("expected no match for nil timestamp")
	}
}

func TestDateFilter_Match_FieldHit(t *testing.T) {
	f, _ := NewDateFilter("2024-03-15", "time")
	e := dateEntry(nil, map[string]any{"time": "2024-03-15T08:00:00Z"})
	if !f.Match(e) {
		t.Fatal("expected match via field")
	}
}

func TestDateFilter_Match_FieldMissing(t *testing.T) {
	f, _ := NewDateFilter("2024-03-15", "time")
	if f.Match(dateEntry(nil, map[string]any{})) {
		t.Fatal("expected no match for missing field")
	}
}
