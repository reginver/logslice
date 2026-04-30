package filter

import (
	"testing"
)

func hostnameEntry(field, value string) map[string]interface{} {
	return map[string]interface{}{field: value}
}

func TestNewHostnameFilter_Valid(t *testing.T) {
	f, err := NewHostnameFilter("host", []string{"web-*", "db-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewHostnameFilter_EmptyField(t *testing.T) {
	_, err := NewHostnameFilter("", []string{"web-*"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewHostnameFilter_NoPatterns(t *testing.T) {
	_, err := NewHostnameFilter("host", []string{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNewHostnameFilter_EmptyPattern(t *testing.T) {
	_, err := NewHostnameFilter("host", []string{"web-*", ""})
	if err == nil {
		t.Fatal("expected error for empty pattern element")
	}
}

func TestHostnameFilter_Match_ExactHit(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"db-01"})
	if !f.Match(hostnameEntry("host", "db-01")) {
		t.Error("expected match for exact pattern")
	}
}

func TestHostnameFilter_Match_GlobHit(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"web-*"})
	if !f.Match(hostnameEntry("host", "web-prod-1")) {
		t.Error("expected match for glob pattern")
	}
}

func TestHostnameFilter_Match_Miss(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"web-*"})
	if f.Match(hostnameEntry("host", "db-01")) {
		t.Error("expected no match")
	}
}

func TestHostnameFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"web-*"})
	if f.Match(map[string]interface{}{"other": "web-01"}) {
		t.Error("expected no match when field absent")
	}
}

func TestHostnameFilter_Match_NonStringValue(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"web-*"})
	if f.Match(map[string]interface{}{"host": 42}) {
		t.Error("expected no match for non-string value")
	}
}

func TestHostnameFilter_Match_SuffixGlob(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"*-prod"})
	if !f.Match(hostnameEntry("host", "web-prod")) {
		t.Error("expected match for suffix glob")
	}
	if f.Match(hostnameEntry("host", "web-staging")) {
		t.Error("expected no match for non-suffix")
	}
}
