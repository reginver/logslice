package filter

import (
	"testing"
)

func cidrGroupEntry(field, ip string) map[string]interface{} {
	if ip == "" {
		return map[string]interface{}{}
	}
	return map[string]interface{}{field: ip}
}

func TestNewCIDRGroupFilter_Valid(t *testing.T) {
	_, err := NewCIDRGroupFilter("src", []string{"10.0.0.0/8", "192.168.0.0/16"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCIDRGroupFilter_EmptyField(t *testing.T) {
	_, err := NewCIDRGroupFilter("", []string{"10.0.0.0/8"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewCIDRGroupFilter_NoCIDRs(t *testing.T) {
	_, err := NewCIDRGroupFilter("src", []string{})
	if err == nil {
		t.Fatal("expected error for empty CIDR list")
	}
}

func TestNewCIDRGroupFilter_InvalidCIDR(t *testing.T) {
	_, err := NewCIDRGroupFilter("src", []string{"not-a-cidr"})
	if err == nil {
		t.Fatal("expected error for invalid CIDR")
	}
}

func TestNewCIDRGroupFilter_EmptyCIDRString(t *testing.T) {
	_, err := NewCIDRGroupFilter("src", []string{""})
	if err == nil {
		t.Fatal("expected error for empty CIDR string")
	}
}

func TestCIDRGroupFilter_Match_Hit(t *testing.T) {
	f, _ := NewCIDRGroupFilter("src", []string{"10.0.0.0/8", "192.168.0.0/16"})
	if !f.Match(cidrGroupEntry("src", "10.1.2.3")) {
		t.Error("expected match for 10.1.2.3 in 10.0.0.0/8")
	}
	if !f.Match(cidrGroupEntry("src", "192.168.5.10")) {
		t.Error("expected match for 192.168.5.10 in 192.168.0.0/16")
	}
}

func TestCIDRGroupFilter_Match_Miss(t *testing.T) {
	f, _ := NewCIDRGroupFilter("src", []string{"10.0.0.0/8"})
	if f.Match(cidrGroupEntry("src", "172.16.0.1")) {
		t.Error("expected no match for 172.16.0.1")
	}
}

func TestCIDRGroupFilter_Match_MissingField(t *testing.T) {
	f, _ := NewCIDRGroupFilter("src", []string{"10.0.0.0/8"})
	if f.Match(cidrGroupEntry("src", "")) {
		t.Error("expected no match when field is absent")
	}
}

func TestCIDRGroupFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewCIDRGroupFilter("src", []string{"10.0.0.0/8"})
	entry := map[string]interface{}{"src": 12345}
	if f.Match(entry) {
		t.Error("expected no match for non-string field value")
	}
}

func TestCIDRGroupFilter_Match_InvalidIP(t *testing.T) {
	f, _ := NewCIDRGroupFilter("src", []string{"10.0.0.0/8"})
	entry := map[string]interface{}{"src": "not-an-ip"}
	if f.Match(entry) {
		t.Error("expected no match for invalid IP string")
	}
}
