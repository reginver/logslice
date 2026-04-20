package filter

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func hashEntry(field, value string) map[string]interface{} {
	return map[string]interface{}{field: value}
}

func md5Prefix(s string, n int) string {
	sum := md5.Sum([]byte(s))
	full := fmt.Sprintf("%x", sum)
	if n > len(full) {
		return full
	}
	return full[:n]
}

func TestNewHashFilter_Valid(t *testing.T) {
	f, err := NewHashFilter("request_id", "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewHashFilter_EmptyField(t *testing.T) {
	_, err := NewHashFilter("", "a1")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewHashFilter_EmptyPrefix(t *testing.T) {
	_, err := NewHashFilter("request_id", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewHashFilter_PrefixTooLong(t *testing.T) {
	_, err := NewHashFilter("request_id", "aabbccddeeff00112233445566778899x")
	if err == nil {
		t.Fatal("expected error for prefix longer than 32 chars")
	}
}

func TestHashFilter_Match_Hit(t *testing.T) {
	value := "user-42"
	prefix := md5Prefix(value, 4)
	f, err := NewHashFilter("user", prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(hashEntry("user", value)) {
		t.Error("expected match")
	}
}

func TestHashFilter_Match_Miss(t *testing.T) {
	f, err := NewHashFilter("user", "0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "user-42" is very unlikely to hash to 0000
	value := "user-42"
	prefix := md5Prefix(value, 4)
	if prefix == "0000" {
		t.Skip("unlucky collision, skipping")
	}
	if f.Match(hashEntry("user", value)) {
		t.Error("expected no match")
	}
}

func TestHashFilter_Match_FieldAbsent(t *testing.T) {
	f, err := NewHashFilter("user", "ab")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Match(map[string]interface{}{"other": "value"}) {
		t.Error("expected no match when field absent")
	}
}

func TestHashFilter_Match_NonStringField(t *testing.T) {
	value := 12345
	prefix := md5Prefix(fmt.Sprintf("%v", value), 4)
	f, err := NewHashFilter("code", prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(map[string]interface{}{"code": value}) {
		t.Error("expected match for numeric field coerced to string")
	}
}
