package cli

import (
	"testing"
)

func TestParseCIDRGroupPairs_Valid(t *testing.T) {
	filters, err := parseCIDRGroupPairs([]string{
		"src=10.0.0.0/8,192.168.0.0/16",
		"dst=172.16.0.0/12",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
}

func TestParseCIDRGroupPairs_Empty(t *testing.T) {
	filters, err := parseCIDRGroupPairs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Fatalf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseCIDRGroupPair_MissingEquals(t *testing.T) {
	_, err := parseCIDRGroupPair("src10.0.0.0/8")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseCIDRGroupPair_EmptyField(t *testing.T) {
	_, err := parseCIDRGroupPair("=10.0.0.0/8")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseCIDRGroupPair_EmptyCIDRList(t *testing.T) {
	_, err := parseCIDRGroupPair("src=")
	if err == nil {
		t.Fatal("expected error for empty CIDR list")
	}
}

func TestParseCIDRGroupPair_InvalidCIDR(t *testing.T) {
	_, err := parseCIDRGroupPair("src=10.0.0.0/8,bad-cidr")
	if err == nil {
		t.Fatal("expected error for invalid CIDR")
	}
}

func TestParseCIDRGroupPair_SingleCIDR(t *testing.T) {
	f, err := parseCIDRGroupPair("src=10.0.0.0/8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
