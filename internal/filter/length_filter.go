package filter

import (
	"fmt"
	"github.com/logslice/logslice/internal/parser"
)

// LengthFilter matches entries where the string length of a field value
// falls within [min, max] (inclusive). Either bound may be -1 to indicate
// no limit.
type LengthFilter struct {
	field string
	min   int
	max   int
}

// NewLengthFilter returns a LengthFilter for the given field and bounds.
// min and max are inclusive; pass -1 to omit a bound.
func NewLengthFilter(field string, min, max int) (*LengthFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("length filter: field name must not be empty")
	}
	if min < -1 {
		return nil, fmt.Errorf("length filter: min must be >= 0 or -1 (unbounded)")
	}
	if max < -1 {
		return nil, fmt.Errorf("length filter: max must be >= 0 or -1 (unbounded)")
	}
	if min != -1 && max != -1 && min > max {
		return nil, fmt.Errorf("length filter: min %d is greater than max %d", min, max)
	}
	return &LengthFilter{field: field, min: min, max: max}, nil
}

// Match returns true when the string representation of the field value has
// a length within the configured bounds.
func (f *LengthFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	l := len(fmt.Sprintf("%v", v))
	if f.min != -1 && l < f.min {
		return false
	}
	if f.max != -1 && l > f.max {
		return false
	}
	return true
}
