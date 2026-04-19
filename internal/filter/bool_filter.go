package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// BoolFilter matches log entries where a field equals a specific boolean value.
type BoolFilter struct {
	field string
	want  bool
}

// NewBoolFilter returns a filter that matches entries where field equals want.
func NewBoolFilter(field string, want bool) (*BoolFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("bool filter: field name must not be empty")
	}
	return &BoolFilter{field: field, want: want}, nil
}

func (f *BoolFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val == f.want
	case string:
		if f.want {
			return val == "true" || val == "1" || val == "yes"
		}
		return val == "false" || val == "0" || val == "no"
	case float64:
		if f.want {
			return val != 0
		}
		return val == 0
	}
	return false
}
