package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// ExistsFilter matches log entries where a given field is present (or absent).
type ExistsFilter struct {
	field  string
	negate bool
}

// NewExistsFilter returns a filter that matches entries containing field.
// If negate is true, it matches entries where the field is absent.
func NewExistsFilter(field string, negate bool) (*ExistsFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("exists filter: field must not be empty")
	}
	return &ExistsFilter{field: field, negate: negate}, nil
}

// Match returns true if the entry satisfies the existence condition.
func (f *ExistsFilter) Match(e parser.LogEntry) bool {
	_, ok := e.Fields[f.field]
	if f.negate {
		return !ok
	}
	return ok
}
