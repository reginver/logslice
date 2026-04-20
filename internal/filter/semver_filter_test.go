package filter

import (
	"testing"

	"github.com/andrewstuart/logslice/internal/parser"
)

func semverEntry(field, value string) parser.LogEntry {
	return parser.LogEntry{Fields: map[string]any{field: value}}
}

func TestNewSemverFilter_Valid(t *testing.T) {
	_, err := NewSemverFilter("version", "1.2.0", "2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSemverFilter_EmptyField(t *testing.T) {
	_, err := NewSemverFilter("", "1.0.0", "2.0.0")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSemverFilter_BothBoundsEmpty(t *testing.T) {
	_, err := NewSemverFilter("version", "", "")
	if err == nil {
		t.Fatal("expected error when both bounds are empty")
	}
}

func TestNewSemverFilter_MaxLessThanMin(t *testing.T) {
	_, err := NewSemverFilter("version", "2.0.0", "1.0.0")
	if err == nil {
		t.Fatal("expected error when max < min")
	}
}

func TestNewSemverFilter_InvalidMin(t *testing.T) {
	_, err := NewSemverFilter("version", "not-a-version", "")
	if err == nil {
		t.Fatal("expected error for invalid min semver")
	}
}

func TestSemverFilter_Match_Hit(t *testing.T) {
	f, _ := NewSemverFilter("version", "1.0.0", "2.0.0")
	if !f.Match(semverEntry("version", "1.5.3")) {
		t.Error("expected match for version within range")
	}
}

func TestSemverFilter_Match_ExactBounds(t *testing.T) {
	f, _ := NewSemverFilter("version", "1.0.0", "2.0.0")
	if !f.Match(semverEntry("version", "1.0.0")) {
		t.Error("expected match for version equal to min")
	}
	if !f.Match(semverEntry("version", "2.0.0")) {
		t.Error("expected match for version equal to max")
	}
}

func TestSemverFilter_Match_Miss_Below(t *testing.T) {
	f, _ := NewSemverFilter("version", "1.2.0", "")
	if f.Match(semverEntry("version", "1.1.9")) {
		t.Error("expected no match for version below min")
	}
}

func TestSemverFilter_Match_Miss_Above(t *testing.T) {
	f, _ := NewSemverFilter("version", "", "1.9.9")
	if f.Match(semverEntry("version", "2.0.0")) {
		t.Error("expected no match for version above max")
	}
}

func TestSemverFilter_Match_FieldAbsent(t *testing.T) {
	f, _ := NewSemverFilter("version", "1.0.0", "2.0.0")
	if f.Match(parser.LogEntry{Fields: map[string]any{"other": "1.5.0"}}) {
		t.Error("expected no match when field is absent")
	}
}

func TestSemverFilter_Match_VPrefix(t *testing.T) {
	f, _ := NewSemverFilter("version", "v1.0.0", "v2.0.0")
	if !f.Match(semverEntry("version", "v1.8.0")) {
		t.Error("expected match for v-prefixed version within range")
	}
}

func TestSemverFilter_Match_NonStringField(t *testing.T) {
	f, _ := NewSemverFilter("version", "1.0.0", "2.0.0")
	e := parser.LogEntry{Fields: map[string]any{"version": 150}}
	if f.Match(e) {
		t.Error("expected no match for non-string field value")
	}
}
