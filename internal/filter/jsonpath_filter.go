package filter

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// JSONPathFilter matches log entries where a dot-separated nested field path
// resolves to a value equal to the expected string.
type JSONPathFilter struct {
	path     []string
	expected string
}

// NewJSONPathFilter creates a filter that checks a nested field path (e.g.
// "request.headers.content-type") against an expected value.
func NewJSONPathFilter(path, expected string) (*JSONPathFilter, error) {
	if path == "" {
		return nil, fmt.Errorf("jsonpath filter: path must not be empty")
	}
	if expected == "" {
		return nil, fmt.Errorf("jsonpath filter: expected value must not be empty")
	}
	parts := strings.Split(path, ".")
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf("jsonpath filter: path segment must not be empty in %q", path)
		}
	}
	return &JSONPathFilter{path: parts, expected: expected}, nil
}

// Match returns true when the nested value at the configured path equals the
// expected string.
func (f *JSONPathFilter) Match(entry parser.LogEntry) bool {
	var current interface{} = entry.Fields
	for _, key := range f.path {
		m, ok := current.(map[string]interface{})
		if !ok {
			return false
		}
		current, ok = m[key]
		if !ok {
			return false
		}
	}
	switch v := current.(type) {
	case string:
		return v == f.expected
	case fmt.Stringer:
		return v.String() == f.expected
	default:
		return fmt.Sprintf("%v", current) == f.expected
	}
}
