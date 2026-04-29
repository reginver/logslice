package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func tenantEntry(field, value string) parser.LogEntry {
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if field != "" {
		e.Fields[field] = value
	}
	return e
}

func TestNewTenantFilter_Valid(t *testing.T) {
	f, err := filter.NewTenantFilter("tenant", []string{"acme", "globex"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewTenantFilter_EmptyField(t *testing.T) {
	_, err := filter.NewTenantFilter("", []string{"acme"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTenantFilter_NoTenants(t *testing.T) {
	_, err := filter.NewTenantFilter("tenant", []string{})
	if err == nil {
		t.Fatal("expected error for empty tenant list")
	}
}

func TestNewTenantFilter_EmptyTenantID(t *testing.T) {
	_, err := filter.NewTenantFilter("tenant", []string{"acme", ""})
	if err == nil {
		t.Fatal("expected error for empty tenant ID")
	}
}

func TestTenantFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewTenantFilter("tenant", []string{"acme", "globex"})
	if !f.Match(tenantEntry("tenant", "acme")) {
		t.Error("expected match for 'acme'")
	}
}

func TestTenantFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := filter.NewTenantFilter("tenant", []string{"Acme"})
	if !f.Match(tenantEntry("tenant", "ACME")) {
		t.Error("expected case-insensitive match")
	}
}

func TestTenantFilter_Match_Miss(t *testing.T) {
	f, _ := filter.NewTenantFilter("tenant", []string{"acme"})
	if f.Match(tenantEntry("tenant", "initech")) {
		t.Error("expected no match for unknown tenant")
	}
}

func TestTenantFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewTenantFilter("tenant", []string{"acme"})
	e := parser.LogEntry{Fields: map[string]interface{}{}}
	if f.Match(e) {
		t.Error("expected no match when field is absent")
	}
}

func TestTenantFilter_Match_NilValue(t *testing.T) {
	f, _ := filter.NewTenantFilter("tenant", []string{"acme"})
	e := parser.LogEntry{Fields: map[string]interface{}{"tenant": nil}}
	if f.Match(e) {
		t.Error("expected no match for nil field value")
	}
}

func TestTenantFilter_Match_NonStringValue(t *testing.T) {
	f, _ := filter.NewTenantFilter("tenant", []string{"acme"})
	e := parser.LogEntry{Fields: map[string]interface{}{"tenant": 42}}
	if f.Match(e) {
		t.Error("expected no match for non-string field value")
	}
}
