package filter

import (
	"fmt"
	"path"

	"github.com/user/logslice/internal/parser"
)

// WildcardFilter matches log entries where a field value matches a glob pattern.
type WildcardFilter struct {
	field   string
	pattern string
}

// NewWildcardFilter creates a WildcardFilter for the given field and glob pattern.
func NewWildcardFilter(field, pattern string) (*WildcardFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("wildcard filter: field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("wildcard filter: pattern must not be empty")
	}
	// Validate pattern syntax.
	if _, err := path.Match(pattern, ""); err != nil {
		return nil, fmt.Errorf("wildcard filter: invalid pattern %q: %w", pattern, err)
	}
	return &WildcardFilter{field: field, pattern: pattern}, nil
}

// Match returns true if the entry's field value matches the glob pattern.
func (f *WildcardFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}
	matched, err := path.Match(f.pattern, s)
	if err != nil {
		return false
	}
	return matched
}
