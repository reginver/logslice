package filter

import (
	"fmt"
	"path"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// PathFilter matches log entries where a field value matches a glob-style path pattern.
type PathFilter struct {
	field   string
	pattern string
}

// NewPathFilter returns a PathFilter that matches entries where the given field
// value matches the provided glob pattern (e.g. "/api/v1/*").
func NewPathFilter(field, pattern string) (*PathFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("path filter: field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("path filter: pattern must not be empty")
	}
	// Validate the pattern is syntactically valid.
	if _, err := path.Match(pattern, ""); err != nil {
		return nil, fmt.Errorf("path filter: invalid pattern %q: %w", pattern, err)
	}
	return &PathFilter{field: field, pattern: pattern}, nil
}

// Match returns true if the entry's field value matches the glob pattern.
func (f *PathFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}
	// Normalise separators to forward slash.
	s = strings.ReplaceAll(s, "\\", "/")
	matched, err := path.Match(f.pattern, s)
	if err != nil {
		return false
	}
	return matched
}
