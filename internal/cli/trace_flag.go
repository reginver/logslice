package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// parseTracePairs parses a slice of "field=spec" strings into TraceFilters.
// spec is a comma-separated list of exact trace IDs or prefix patterns (e.g.
// "abc123,def*").
func parseTracePairs(pairs []string) ([]*filter.TraceFilter, error) {
	var filters []*filter.TraceFilter
	for _, p := range pairs {
		f, err := parseTracePair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseTracePair(pair string) (*filter.TraceFilter, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return nil, fmt.Errorf("trace filter: expected field=spec, got %q", pair)
	}
	field := pair[:idx]
	spec := pair[idx+1:]
	if field == "" {
		return nil, fmt.Errorf("trace filter: field must not be empty in %q", pair)
	}
	if spec == "" {
		return nil, fmt.Errorf("trace filter: spec must not be empty in %q", pair)
	}
	return filter.NewTraceFilter(field, spec)
}
