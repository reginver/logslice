package filter

import (
	"testing"
)

func uaEntry(field, value string) map[string]interface{} {
	return map[string]interface{}{field: value}
}

func TestNewUserAgentFilter_Valid(t *testing.T) {
	f, err := NewUserAgentFilter("ua", []string{"mozilla", "chrome"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewUserAgentFilter_EmptyField(t *testing.T) {
	_, err := NewUserAgentFilter("", []string{"chrome"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewUserAgentFilter_NoTokens(t *testing.T) {
	_, err := NewUserAgentFilter("ua", []string{})
	if err == nil {
		t.Fatal("expected error for empty token list")
	}
}

func TestNewUserAgentFilter_EmptyToken(t *testing.T) {
	_, err := NewUserAgentFilter("ua", []string{"chrome", ""})
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestUserAgentFilter_Match_Hit(t *testing.T) {
	f, _ := NewUserAgentFilter("ua", []string{"chrome"})
	entry := uaEntry("ua", "Mozilla/5.0 (Windows) AppleWebKit Chrome/114")
	if !f.Match(entry) {
		t.Fatal("expected match")
	}
}

func TestUserAgentFilter_Match_Miss(t *testing.T) {
	f, _ := NewUserAgentFilter("ua", []string{"chrome"})
	entry := uaEntry("ua", "curl/7.88.1")
	if f.Match(entry) {
		t.Fatal("expected no match")
	}
}

func TestUserAgentFilter_Match_CaseInsensitive(t *testing.T) {
	f, _ := NewUserAgentFilter("ua", []string{"FIREFOX"})
	entry := uaEntry("ua", "Mozilla/5.0 Gecko Firefox/115")
	if !f.Match(entry) {
		t.Fatal("expected case-insensitive match")
	}
}

func TestUserAgentFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewUserAgentFilter("ua", []string{"chrome"})
	if f.Match(map[string]interface{}{"other": "value"}) {
		t.Fatal("expected no match when field absent")
	}
}

func TestUserAgentFilter_Match_NonStringValue(t *testing.T) {
	f, _ := NewUserAgentFilter("ua", []string{"chrome"})
	if f.Match(map[string]interface{}{"ua": 42}) {
		t.Fatal("expected no match for non-string value")
	}
}
