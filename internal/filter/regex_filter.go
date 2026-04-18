package filter

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// RegexFilter matches log entries where a named field matches a regular expression.
type RegexFilter struct {
	field   string
	pattern *regexp.Regexp
}

// NewRegexFilter compiles the given pattern and returns a RegexFilter.
// Returns an error if the field is empty or the pattern is invalid.
func NewRegexFilter(field, pattern string) (*RegexFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("regex filter: field name must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("regex filter: invalid pattern %q: %w", pattern, err)
	}
	return &RegexFilter{field: field, pattern: re}, nil
}

// Match returns true when the entry's field value matches the compiled pattern.
func (f *RegexFilter) Match(entry parser.LogEntry) bool {
	val, ok := entry.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		s = fmt.Sprintf("%v", val)
	}
	return f.pattern.MatchString(s)
}
