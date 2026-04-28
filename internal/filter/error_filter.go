package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// ErrorFilter matches log entries whose error field contains one of the
// specified error codes or substrings. Matching is case-insensitive.
type ErrorFilter struct {
	field  string
	codes  []string
}

// NewErrorFilter creates a filter that matches entries where the given field
// contains at least one of the provided error codes (case-insensitive).
func NewErrorFilter(field string, codes []string) (*ErrorFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("error filter: field must not be empty")
	}
	if len(codes) == 0 {
		return nil, fmt.Errorf("error filter: at least one error code is required")
	}
	norm := make([]string, 0, len(codes))
	for _, c := range codes {
		c = strings.TrimSpace(c)
		if c == "" {
			return nil, fmt.Errorf("error filter: error code must not be empty")
		}
		norm = append(norm, strings.ToLower(c))
	}
	return &ErrorFilter{field: field, codes: norm}, nil
}

// Match returns true when the entry's field value (lowercased) contains any
// of the registered error codes as a substring.
func (f *ErrorFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	val := strings.ToLower(fmt.Sprintf("%v", v))
	for _, code := range f.codes {
		if strings.Contains(val, code) {
			return true
		}
	}
	return false
}
