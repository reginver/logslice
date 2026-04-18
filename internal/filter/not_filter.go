package filter

import (
	"fmt"

	"github.com/logslice/logslice/internal/parser"
)

// NotFilter negates the result of an inner filter.
type NotFilter struct {
	inner Filter
}

// NewNotFilter returns a Filter that matches entries NOT matched by inner.
func NewNotFilter(inner Filter) (*NotFilter, error) {
	if inner == nil {
		return nil, fmt.Errorf("not_filter: inner filter must not be nil")
	}
	return &NotFilter{inner: inner}, nil
}

// Match returns true if the inner filter does NOT match the entry.
func (f *NotFilter) Match(e parser.LogEntry) bool {
	return !f.inner.Match(e)
}
