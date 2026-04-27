package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// LatencyFilter matches log entries whose numeric field (interpreted as
// milliseconds) falls within [min, max]. Either bound may be zero to indicate
// unbounded on that side.
type LatencyFilter struct {
	field string
	min   time.Duration
	max   time.Duration
}

// NewLatencyFilter creates a LatencyFilter. min and max are millisecond
// durations. Pass 0 for min to skip the lower bound; pass 0 for max to skip
// the upper bound. At least one bound must be non-zero.
func NewLatencyFilter(field string, min, max time.Duration) (*LatencyFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("latency filter: field name must not be empty")
	}
	if min == 0 && max == 0 {
		return nil, fmt.Errorf("latency filter: at least one of min or max must be set")
	}
	if min < 0 {
		return nil, fmt.Errorf("latency filter: min must not be negative")
	}
	if max < 0 {
		return nil, fmt.Errorf("latency filter: max must not be negative")
	}
	if min > 0 && max > 0 && min > max {
		return nil, fmt.Errorf("latency filter: min (%v) must not exceed max (%v)", min, max)
	}
	return &LatencyFilter{field: field, min: min, max: max}, nil
}

// Match returns true when the entry's field value (treated as milliseconds)
// falls within the configured range.
func (f *LatencyFilter) Match(e parser.Entry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	var ms float64
	switch val := v.(type) {
	case float64:
		ms = val
	case int:
		ms = float64(val)
	case int64:
		ms = float64(val)
	default:
		return false
	}
	d := time.Duration(ms) * time.Millisecond
	if f.min > 0 && d < f.min {
		return false
	}
	if f.max > 0 && d > f.max {
		return false
	}
	return true
}
