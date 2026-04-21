package filter

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func jpEntry(fields map[string]interface{}) parser.LogEntry {
	return parser.LogEntry{
		Timestamp: func() *time.Time { t := time.Now(); return &t }(),
		Fields:    fields,
	}
}

func TestNewJSONPathFilter_Valid(t *testing.T) {
	f, err := NewJSONPathFilter("request.method", "GET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewJSONPathFilter_EmptyPath(t *testing.T) {
	_, err := NewJSONPathFilter("", "GET")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewJSONPathFilter_EmptyExpected(t *testing.T) {
	_, err := NewJSONPathFilter("request.method", "")
	if err == nil {
		t.Fatal("expected error for empty expected value")
	}
}

func TestNewJSONPathFilter_EmptySegment(t *testing.T) {
	_, err := NewJSONPathFilter("request..method", "GET")
	if err == nil {
		t.Fatal("expected error for empty path segment")
	}
}

func TestJSONPathFilter_Match_Hit(t *testing.T) {
	f, _ := NewJSONPathFilter("request.method", "GET")
	entry := jpEntry(map[string]interface{}{
		"request": map[string]interface{}{"method": "GET"},
	})
	if !f.Match(entry) {
		t.Error("expected match")
	}
}

func TestJSONPathFilter_Match_Miss(t *testing.T) {
	f, _ := NewJSONPathFilter("request.method", "POST")
	entry := jpEntry(map[string]interface{}{
		"request": map[string]interface{}{"method": "GET"},
	})
	if f.Match(entry) {
		t.Error("expected no match")
	}
}

func TestJSONPathFilter_Match_MissingIntermediate(t *testing.T) {
	f, _ := NewJSONPathFilter("request.method", "GET")
	entry := jpEntry(map[string]interface{}{"level": "info"})
	if f.Match(entry) {
		t.Error("expected no match when intermediate key is absent")
	}
}

func TestJSONPathFilter_Match_NonMapIntermediate(t *testing.T) {
	f, _ := NewJSONPathFilter("request.method", "GET")
	entry := jpEntry(map[string]interface{}{
		"request": "not-a-map",
	})
	if f.Match(entry) {
		t.Error("expected no match when intermediate value is not a map")
	}
}

func TestJSONPathFilter_Match_DeepPath(t *testing.T) {
	f, _ := NewJSONPathFilter("a.b.c", "deep")
	entry := jpEntry(map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{"c": "deep"},
		},
	})
	if !f.Match(entry) {
		t.Error("expected match on deep path")
	}
}
