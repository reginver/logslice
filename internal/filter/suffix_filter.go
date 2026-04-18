package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// SuffixFilter matches log entries where a field value ends with a given suffix.
type SuffixFilter struct {
	field  string
	suffix string
}

// NewSuffixFilter creates a SuffixFilter for the given field and suffix.
func NewSuffixFilter(field, suffix string) (*SuffixFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("suffix filter: field must not be empty")
	}
	if suffix == "" {
		return nil, fmt.Errorf("suffix filter: suffix must not be empty")
	}
	return &SuffixFilter{field: field, suffix: suffix}, nil
}

// Match returns true if the entry's field value ends with the configured suffix.
func (f *SuffixFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}
	return strings.HasSuffix(s, f.suffix)
}
