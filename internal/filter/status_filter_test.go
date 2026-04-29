package filter

import (
	"testing"
)

func statusEntry(field string, value interface{}) map[string]interface{} {
	return map[string]interface{}{field: value}
}

func TestNewStatusFilter_Valid(t *testing.T) {
	f, err := NewStatusFilter("status", []string{"200-299", "404"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewStatusFilter_EmptyField(t *testing.T) {
	_, err := NewStatusFilter("", []string{"200"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewStatusFilter_NoSpecs(t *testing.T) {
	_, err := NewStatusFilter("status", []string{})
	if err == nil {
		t.Fatal("expected error for empty specs")
	}
}

func TestNewStatusFilter_EmptySpec(t *testing.T) {
	_, err := NewStatusFilter("status", []string{""})
	if err == nil {
		t.Fatal("expected error for blank spec")
	}
}

func TestNewStatusFilter_InvalidCode(t *testing.T) {
	_, err := NewStatusFilter("status", []string{"abc"})
	if err == nil {
		t.Fatal("expected error for non-numeric code")
	}
}

func TestNewStatusFilter_LoGreaterThanHi(t *testing.T) {
	_, err := NewStatusFilter("status", []string{"500-200"})
	if err == nil {
		t.Fatal("expected error when lo > hi")
	}
}

func TestStatusFilter_Match_Float64(t *testing.T) {
	f, _ := NewStatusFilter("status", []string{"200-299"})
	if !f.Match(statusEntry("status", float64(200))) {
		t.Error("expected match for 200")
	}
	if !f.Match(statusEntry("status", float64(299))) {
		t.Error("expected match for 299")
	}
	if f.Match(statusEntry("status", float64(300))) {
		t.Error("expected no match for 300")
	}
}

func TestStatusFilter_Match_String(t *testing.T) {
	f, _ := NewStatusFilter("code", []string{"404"})
	if !f.Match(statusEntry("code", "404")) {
		t.Error("expected match for string '404'")
	}
	if f.Match(statusEntry("code", "200")) {
		t.Error("expected no match for string '200'")
	}
}

func TestStatusFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewStatusFilter("status", []string{"200"})
	if f.Match(map[string]interface{}{"other": 200}) {
		t.Error("expected no match when field absent")
	}
}

func TestStatusFilter_Match_MultiRange(t *testing.T) {
	f, _ := NewStatusFilter("status", []string{"200-299", "400-499"})
	if !f.Match(statusEntry("status", float64(201))) {
		t.Error("expected match in first range")
	}
	if !f.Match(statusEntry("status", float64(404))) {
		t.Error("expected match in second range")
	}
	if f.Match(statusEntry("status", float64(301))) {
		t.Error("expected no match outside both ranges")
	}
}
