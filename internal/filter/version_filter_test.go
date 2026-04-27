package filter

import (
	"testing"
)

func versionEntry(field, val string) map[string]interface{} {
	return map[string]interface{}{field: val}
}

func TestNewVersionFilter_Valid(t *testing.T) {
	_, err := NewVersionFilter("version", "1.0.0", "2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewVersionFilter_EmptyField(t *testing.T) {
	_, err := NewVersionFilter("", "1.0.0", "2.0.0")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewVersionFilter_BothBoundsEmpty(t *testing.T) {
	_, err := NewVersionFilter("version", "", "")
	if err == nil {
		t.Fatal("expected error when both bounds are empty")
	}
}

func TestNewVersionFilter_InvalidMin(t *testing.T) {
	_, err := NewVersionFilter("version", "abc", "2.0.0")
	if err == nil {
		t.Fatal("expected error for invalid min")
	}
}

func TestNewVersionFilter_MinGreaterThanMax(t *testing.T) {
	_, err := NewVersionFilter("version", "3.0.0", "1.0.0")
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestVersionFilter_Match_Hit(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.2.0", "2.0.0")
	if !f.Match(versionEntry("version", "1.5.3")) {
		t.Error("expected match for version within range")
	}
}

func TestVersionFilter_Match_BelowMin(t *testing.T) {
	f, _ := NewVersionFilter("version", "2.0.0", "3.0.0")
	if f.Match(versionEntry("version", "1.9.9")) {
		t.Error("expected no match for version below min")
	}
}

func TestVersionFilter_Match_AboveMax(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.0.0", "1.9.9")
	if f.Match(versionEntry("version", "2.0.0")) {
		t.Error("expected no match for version above max")
	}
}

func TestVersionFilter_Match_ExactBounds(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.0.0", "1.0.0")
	if !f.Match(versionEntry("version", "1.0.0")) {
		t.Error("expected match for exact boundary version")
	}
}

func TestVersionFilter_Match_MissingField(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.0.0", "2.0.0")
	if f.Match(map[string]interface{}{"other": "1.5.0"}) {
		t.Error("expected no match when field is absent")
	}
}

func TestVersionFilter_Match_InvalidFieldValue(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.0.0", "2.0.0")
	if f.Match(versionEntry("version", "not-a-version")) {
		t.Error("expected no match for non-version string")
	}
}

func TestVersionFilter_UnboundedMin(t *testing.T) {
	f, _ := NewVersionFilter("version", "", "2.0.0")
	if !f.Match(versionEntry("version", "0.1.0")) {
		t.Error("expected match when min is unbounded")
	}
}

func TestVersionFilter_UnboundedMax(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.0.0", "")
	if !f.Match(versionEntry("version", "99.0.0")) {
		t.Error("expected match when max is unbounded")
	}
}
