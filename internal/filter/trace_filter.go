package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// TraceFilter matches log entries whose trace ID field has a given prefix or
// matches one of a set of exact trace IDs.
type TraceFilter struct {
	field   string
	prefixes []string
	exact    map[string]struct{}
}

// NewTraceFilter creates a TraceFilter. spec is a comma-separated list of
// trace ID prefixes or exact IDs to match against the value of field.
func NewTraceFilter(field, spec string) (*TraceFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("trace filter: field must not be empty")
	}
	if spec == "" {
		return nil, fmt.Errorf("trace filter: spec must not be empty")
	}
	parts := strings.Split(spec, ",")
	var prefixes []string
	exact := make(map[string]struct{})
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("trace filter: empty trace ID in spec")
		}
		if strings.HasSuffix(p, "*") {
			prefixes = append(prefixes, strings.TrimSuffix(p, "*"))
		} else {
			exact[p] = struct{}{}
		}
	}
	if len(prefixes) == 0 && len(exact) == 0 {
		return nil, fmt.Errorf("trace filter: no valid trace IDs provided")
	}
	return &TraceFilter{field: field, prefixes: prefixes, exact: exact}, nil
}

// Match returns true if the entry's trace field value matches any exact ID or
// prefix defined in the filter.
func (f *TraceFilter) Match(e parser.Entry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	if _, found := f.exact[s]; found {
		return true
	}
	for _, pfx := range f.prefixes {
		if strings.HasPrefix(s, pfx) {
			return true
		}
	}
	return false
}
