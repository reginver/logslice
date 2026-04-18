package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// PrefixFilter matches log entries where a string field starts with a given prefix.
type PrefixFilter struct {
	field  string
	prefix string
}

// NewPrefixFilter creates a PrefixFilter for the given field and prefix.
// Returns an error if field or prefix is empty.
func NewPrefixFilter(field, prefix string) (*PrefixFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("prefix filter: field must not be empty")
	}
	if prefix == "" {
		return nil, fmt.Errorf("prefix filter: prefix must not be empty")
	}
	return &PrefixFilter{field: field, prefix: prefix}, nil
}

// Match returns true if the entry's field value starts with the configured prefix.
func (f *PrefixFilter) Match(e parser.LogEntry) bool {
	val, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	str, ok := val.(string)
	if !ok {
		return false
	}
	return strings.HasPrefix(str, f.prefix)
}
