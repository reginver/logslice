package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parsePathPairs parses a slice of "field=pattern" strings into PathFilters.
// Example: ["url=/api/v1/*", "path=/static/**"]
func parsePathPairs(pairs []string) ([]*filter.PathFilter, error) {
	var filters []*filter.PathFilter
	for _, p := range pairs {
		f, err := parsePathPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parsePathPair(pair string) (*filter.PathFilter, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return nil, fmt.Errorf("path flag: expected field=pattern, got %q", pair)
	}
	field := pair[:idx]
	pattern := pair[idx+1:]
	if field == "" {
		return nil, fmt.Errorf("path flag: field must not be empty in %q", pair)
	}
	if pattern == "" {
		return nil, fmt.Errorf("path flag: pattern must not be empty in %q", pair)
	}
	return filter.NewPathFilter(field, pattern)
}
