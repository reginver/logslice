package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func jtEntry(fields map[string]interface{}) parser.LogEntry {
	return parser.LogEntry{Fields: fields}
}

func TestNewJSONTypeFilter_Valid(t *testing.T) {
	f, err := NewJSONTypeFilter("level", "string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewJSONTypeFilter_EmptyField(t *testing.T) {
	_, err := NewJSONTypeFilter("", "string")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewJSONTypeFilter_UnknownType(t *testing.T) {
	_, err := NewJSONTypeFilter("field", "integer")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestJSONTypeFilter_Match_String(t *testing.T) {
	f, _ := NewJSONTypeFilter("msg", "string")
	if !f.Match(jtEntry(map[string]interface{}{"msg": "hello"})) {
		t.Error("expected match for string value")
	}
	if f.Match(jtEntry(map[string]interface{}{"msg": float64(42)})) {
		t.Error("expected no match for number value")
	}
}

func TestJSONTypeFilter_Match_Number(t *testing.T) {
	f, _ := NewJSONTypeFilter("code", "number")
	if !f.Match(jtEntry(map[string]interface{}{"code": float64(200)})) {
		t.Error("expected match for number value")
	}
}

func TestJSONTypeFilter_Match_Bool(t *testing.T) {
	f, _ := NewJSONTypeFilter("ok", "bool")
	if !f.Match(jtEntry(map[string]interface{}{"ok": true})) {
		t.Error("expected match for bool true")
	}
	if !f.Match(jtEntry(map[string]interface{}{"ok": false})) {
		t.Error("expected match for bool false")
	}
}

func TestJSONTypeFilter_Match_Null(t *testing.T) {
	f, _ := NewJSONTypeFilter("data", "null")
	if !f.Match(jtEntry(map[string]interface{}{"data": nil})) {
		t.Error("expected match for nil value")
	}
	// missing field also counts as null
	if !f.Match(jtEntry(map[string]interface{}{})) {
		t.Error("expected match when field is absent")
	}
}

func TestJSONTypeFilter_Match_Array(t *testing.T) {
	f, _ := NewJSONTypeFilter("tags", "array")
	if !f.Match(jtEntry(map[string]interface{}{"tags": []interface{}{"a", "b"}})) {
		t.Error("expected match for array value")
	}
}

func TestJSONTypeFilter_Match_Object(t *testing.T) {
	f, _ := NewJSONTypeFilter("meta", "object")
	if !f.Match(jtEntry(map[string]interface{}{"meta": map[string]interface{}{"k": "v"}})) {
		t.Error("expected match for object value")
	}
}
