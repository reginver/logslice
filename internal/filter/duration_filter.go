package filter

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// DurationFilter matches log entries where a named field, interpreted as a
// duration string (e.g. "1.5s", "300ms"), falls within [min, max].
// Either bound may be zero to indicate unbounded.
type DurationFilter struct {
	field string
	min   time.Duration
	max   time.Duration
}

// NewDurationFilter creates a DurationFilter for the given field and bounds.
// At least one of min or max must be non-zero. If both are set, min <= max.
func NewDurationFilter(field string, min, max time.Duration) (*DurationFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("duration filter: field name must not be empty")
	}
	if min == 0 && max == 0 {
		return nil, fmt.Errorf("duration filter: at least one of min or max must be set")
	}
	if min != 0 && max != 0 && min > max {
		return nil, fmt.Errorf("duration filter: min %v is greater than max %v", min, max)
	}
	return &DurationFilter{field: field, min: min, max: max}, nil
}

// Match returns true when the entry's field value parses as a duration within
// the configured bounds.
func (f *DurationFilter) Match(entry parser.LogEntry) bool {
	v, ok := entry.Fields[f.field]
	if !ok {
		return false
	}
	var raw string
	switch s := v.(type) {
	case string:
		raw = s
	default:
		return false
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return false
	}
	if f.min != 0 && d < f.min {
		return false
	}
	if f.max != 0 && d > f.max {
		return false
	}
	return true
}
