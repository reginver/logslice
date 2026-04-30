package filter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func checksumEntry(field, value string) map[string]interface{} {
	return map[string]interface{}{field: value}
}

func md5Of(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func sha1Of(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func sha256Of(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func TestNewChecksumFilter_Valid(t *testing.T) {
	f, err := NewChecksumFilter("id", "md5", "ab")
	if err != nil || f == nil {
		t.Fatalf("expected valid filter, got err=%v", err)
	}
}

func TestNewChecksumFilter_EmptyField(t *testing.T) {
	_, err := NewChecksumFilter("", "md5", "ab")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewChecksumFilter_UnknownAlgo(t *testing.T) {
	_, err := NewChecksumFilter("id", "crc32", "ab")
	if err == nil {
		t.Fatal("expected error for unknown algorithm")
	}
}

func TestNewChecksumFilter_EmptyPrefix(t *testing.T) {
	_, err := NewChecksumFilter("id", "sha256", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestChecksumFilter_MD5_Hit(t *testing.T) {
	val := "hello"
	digest := md5Of(val)
	f, _ := NewChecksumFilter("msg", "md5", digest[:6])
	if !f.Match(checksumEntry("msg", val)) {
		t.Error("expected match")
	}
}

func TestChecksumFilter_SHA1_Hit(t *testing.T) {
	val := "world"
	digest := sha1Of(val)
	f, _ := NewChecksumFilter("msg", "sha1", digest[:8])
	if !f.Match(checksumEntry("msg", val)) {
		t.Error("expected match")
	}
}

func TestChecksumFilter_SHA256_Miss(t *testing.T) {
	f, _ := NewChecksumFilter("msg", "sha256", "0000000000")
	// highly unlikely real digest starts with ten zeros
	if f.Match(checksumEntry("msg", "logslice")) {
		digest := sha256Of("logslice")
		if len(digest) >= 10 && digest[:10] == "0000000000" {
			t.Skip("unlikely collision occurred")
		}
		t.Error("unexpected match")
	}
}

func TestChecksumFilter_FieldAbsent(t *testing.T) {
	f, _ := NewChecksumFilter("id", "md5", "ab")
	if f.Match(map[string]interface{}{"other": "val"}) {
		t.Error("expected no match when field absent")
	}
}

func TestChecksumFilter_AlgoCaseInsensitive(t *testing.T) {
	val := "test"
	digest := md5Of(val)
	f, err := NewChecksumFilter("f", "MD5", digest[:4])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(checksumEntry("f", val)) {
		t.Error("expected match with uppercase algo name")
	}
}
