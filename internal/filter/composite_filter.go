package filter

import "github.com/yourorg/logslice/internal/parser"

// Filter is the interface implemented by all log entry filters.
type Filter interface {
	Match(entry parser.LogEntry) bool
}

// CompositeFilter combines multiple filters with AND semantics.
type CompositeFilter struct {
	filters []Filter
}

// NewCompositeFilter returns a CompositeFilter that passes entries
// only when all provided filters match.
func NewCompositeFilter(filters ...Filter) *CompositeFilter {
	return &CompositeFilter{filters: filters}
}

// Match returns true if every child filter matches the entry.
func (c *CompositeFilter) Match(entry parser.LogEntry) bool {
	for _, f := range c.filters {
		if !f.Match(entry) {
			return false
		}
	}
	return true
}

// AnyFilter combines multiple filters with OR semantics.
type AnyFilter struct {
	filters []Filter
}

// NewAnyFilter returns an AnyFilter that passes entries when at
// least one child filter matches.
func NewAnyFilter(filters ...Filter) *AnyFilter {
	return &AnyFilter{filters: filters}
}

// Match returns true if any child filter matches the entry.
func (a *AnyFilter) Match(entry parser.LogEntry) bool {
	for _, f := range a.filters {
		if f.Match(entry) {
			return true
		}
	}
	return false
}
