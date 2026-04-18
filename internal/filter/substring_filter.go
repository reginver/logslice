package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// SubstringFilter matches log entries where a field contains a given substring.
type SubstringFilter struct {
	field     string
	substring string
	caseFold  bool
}

// NewSubstringFilter creates a SubstringFilter for the given field and substring.
// If caseFold is true, matching is case-insensitive.
func NewSubstringFilter(field, substring string, caseFold bool) (*SubstringFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("substring filter: field must not be empty")
	}
	if substring == "" {
		return nil, fmt.Errorf("substring filter: substring must not be empty")
	}
	return &SubstringFilter{field: field, substring: substring, caseFold: caseFold}, nil
}

// Match returns true if the entry's field value contains the substring.
func (f *SubstringFilter) Match(e parser.LogEntry) bool {
	val, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		s = fmt.Sprintf("%v", val)
	}
	if f.caseFold {
		return strings.Contains(strings.ToLower(s), strings.ToLower(f.substring))
	}
	return strings.Contains(s, f.substring)
}
