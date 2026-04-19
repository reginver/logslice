package filter_test

import (
	"testing"

	"github.com/your/logslice/internal/filter"
	"github.com/your/logslice/internal/parser"
)

func ipEntry(ip string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{"client_ip": ip}}
}

func TestNewIPFilter_Valid(t *testing.T) {
	_, err := filter.NewIPFilter("client_ip", "10.0.0.0/8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewIPFilter_EmptyField(t *testing.T) {
	_, err := filter.NewIPFilter("", "10.0.0.0/8")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewIPFilter_InvalidCIDR(t *testing.T) {
	_, err := filter.NewIPFilter("client_ip", "not-a-cidr")
	if err == nil {
		t.Fatal("expected error for invalid CIDR")
	}
}

func TestNewIPFilter_EmptyCIDR(t *testing.T) {
	_, err := filter.NewIPFilter("client_ip", "")
	if err == nil {
		t.Fatal("expected error for empty CIDR")
	}
}

func TestIPFilter_Match_Hit(t *testing.T) {
	f, _ := filter.NewIPFilter("client_ip", "192.168.1.0/24")
	if !f.Match(ipEntry("192.168.1.42")) {
		t.Error("expected match")
	}
}

func TestIPFilter_Match_Miss(t *testing.T) {
	f, _ := filter.NewIPFilter("client_ip", "192.168.1.0/24")
	if f.Match(ipEntry("10.0.0.1")) {
		t.Error("expected no match")
	}
}

func TestIPFilter_Match_InvalidIP(t *testing.T) {
	f, _ := filter.NewIPFilter("client_ip", "192.168.1.0/24")
	if f.Match(ipEntry("not-an-ip")) {
		t.Error("expected no match for invalid IP")
	}
}

func TestIPFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := filter.NewIPFilter("client_ip", "192.168.1.0/24")
	e := parser.LogEntry{Fields: map[string]any{}}
	if f.Match(e) {
		t.Error("expected no match when field absent")
	}
}

func TestIPFilter_Match_IPv6(t *testing.T) {
	f, _ := filter.NewIPFilter("client_ip", "2001:db8::/32")
	if !f.Match(ipEntry("2001:db8::1")) {
		t.Error("expected match for IPv6")
	}
}
