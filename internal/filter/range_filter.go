package filter

import (
	"fmt"

	"github.com/example/logslice/internal/parser"
)

// RangeFilter filters log entries where a numeric field falls within [Min, Max].
type RangeFilter struct {
	Field string
	Min   float64
	Max   float64
}

// NewRangeFilter creates a RangeFilter for the given field and numeric bounds.
func NewRangeFilter(field string, min, max float64) (*RangeFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("range filter: field name must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("range filter: min %.4g is greater than max %.4g", min, max)
	}
	return &RangeFilter{Field: field, Min: min, Max: max}, nil
}

// Match returns true when the entry's field value is a number within [Min, Max].
func (f *RangeFilter) Match(e parser.LogEntry) bool {
	raw, ok := e.Fields[f.Field]
	if !ok {
		return false
	}
	var v float64
	switch n := raw.(type) {
	case float64:
		v = n
	case int:
		v = float64(n)
	case int64:
		v = float64(n)
	default:
		return false
	}
	return v >= f.Min && v <= f.Max
}
