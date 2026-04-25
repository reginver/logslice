package filter

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func TestTimeOfDayFilter_WithComposite(t *testing.T) {
	tod, err := NewTimeOfDayFilter("08:00", "12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	field, err := NewFieldFilter("env", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	composite := NewCompositeFilter(tod, field)

	match := func(hour int, env string) bool {
		ts := time.Date(2024, 6, 1, hour, 0, 0, 0, time.UTC)
		e := &parser.Entry{
			Timestamp: &ts,
			Fields:    map[string]interface{}{"env": env},
		}
		return composite.Match(e)
	}

	if !match(10, "prod") {
		t.Error("expected match: 10:00 + prod")
	}
	if match(14, "prod") {
		t.Error("expected no match: 14:00 outside window")
	}
	if match(10, "staging") {
		t.Error("expected no match: wrong env")
	}
}

func TestTimeOfDayFilter_Negated(t *testing.T) {
	tod, err := NewTimeOfDayFilter("09:00", "17:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	not := NewNotFilter(tod)

	offHour := time.Date(2024, 6, 1, 20, 0, 0, 0, time.UTC)
	bizHour := time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC)

	if !not.Match(&parser.Entry{Timestamp: &offHour}) {
		t.Error("expected match for time outside business hours when negated")
	}
	if not.Match(&parser.Entry{Timestamp: &bizHour}) {
		t.Error("expected no match for business hours when negated")
	}
}
