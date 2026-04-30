package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func httpMethodEntry(field, value string) *parser.Entry {
	return &parser.Entry{Fields: map[string]interface{}{field: value}}
}

func TestNewHTTPMethodFilter_Valid(t *testing.T) {
	f, err := NewHTTPMethodFilter("method", []string{"GET", "POST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewHTTPMethodFilter_EmptyField(t *testing.T) {
	_, err := NewHTTPMethodFilter("", []string{"GET"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewHTTPMethodFilter_NoMethods(t *testing.T) {
	_, err := NewHTTPMethodFilter("method", []string{})
	if err == nil {
		t.Fatal("expected error for empty method list")
	}
}

func TestNewHTTPMethodFilter_EmptyMethod(t *testing.T) {
	_, err := NewHTTPMethodFilter("method", []string{"GET", ""})
	if err == nil {
		t.Fatal("expected error for empty method string")
	}
}

func TestHTTPMethodFilter_Match_Hit(t *testing.T) {
	f, _ := NewHTTPMethodFilter("method", []string{"GET", "POST"})
	if !f.Match(httpMethodEntry("method", "GET")) {
		t.Error("expected match for GET")
	}
	if !f.Match(httpMethodEntry("method", "post")) {
		t.Error("expected match for post (case-insensitive)")
	}
}

func TestHTTPMethodFilter_Match_Miss(t *testing.T) {
	f, _ := NewHTTPMethodFilter("method", []string{"GET", "POST"})
	if f.Match(httpMethodEntry("method", "DELETE")) {
		t.Error("expected no match for DELETE")
	}
}

func TestHTTPMethodFilter_Match_MissingField(t *testing.T) {
	f, _ := NewHTTPMethodFilter("method", []string{"GET"})
	e := &parser.Entry{Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestHTTPMethodFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewHTTPMethodFilter("method", []string{"GET"})
	e := &parser.Entry{Fields: map[string]interface{}{"method": 42}}
	if f.Match(e) {
		t.Error("expected no match for non-string field value")
	}
}
