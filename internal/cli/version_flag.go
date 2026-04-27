package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// parseVersionPairs parses a slice of "field=min:max" strings into
// version range filters. Either min or max may be omitted (e.g. "field=:2.0.0"
// or "field=1.0.0:").
func parseVersionPairs(pairs []string) ([]filter.Filter, error) {
	filters := make([]filter.Filter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseVersionPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseVersionPair(pair string) (filter.Filter, error) {
	eqIdx := strings.IndexByte(pair, '=')
	if eqIdx < 0 {
		return nil, fmt.Errorf("version filter: expected field=min:max, got %q", pair)
	}
	field := pair[:eqIdx]
	if field == "" {
		return nil, fmt.Errorf("version filter: field must not be empty in %q", pair)
	}
	rest := pair[eqIdx+1:]
	colIdx := strings.IndexByte(rest, ':')
	if colIdx < 0 {
		return nil, fmt.Errorf("version filter: expected min:max after '=', got %q", pair)
	}
	min := rest[:colIdx]
	max := rest[colIdx+1:]
	if min == "" && max == "" {
		return nil, fmt.Errorf("version filter: at least one of min or max must be set in %q", pair)
	}
	return filter.NewVersionFilter(field, min, max)
}
