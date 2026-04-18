package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// containsFilter matches entries where a field contains all of the given values.
type containsFilter struct {
	field  string
	values []string
}

// NewContainsFilter returns a Filter that matches when the named field's string
// representation contains every value in values (case-insensitive).
func NewContainsFilter(field string, values []string) (parser.Filter, error) {
	if field == "" {
		return nil, fmt.Errorf("contains filter: field name must not be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("contains filter: at least one value required")
	}
	norm := make([]string, len(values))
	for i, v := range values {
		if v == "" {
			return nil, fmt.Errorf("contains filter: value at index %d must not be empty", i)
		}
		norm[i] = strings.ToLower(v)
	}
	return &containsFilter{field: field, values: norm}, nil
}

func (f *containsFilter) Match(e *parser.Entry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s := strings.ToLower(fmt.Sprintf("%v", v))
	for _, needle := range f.values {
		if !strings.Contains(s, needle) {
			return false
		}
	}
	return true
}
