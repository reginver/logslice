package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logslice/logslice/internal/filter"
)

// parseLengthPairs parses a slice of "field:min:max" strings into LengthFilters.
// Use -1 for min or max to indicate an unbounded side.
func parseLengthPairs(pairs []string) ([]*filter.LengthFilter, error) {
	var filters []*filter.LengthFilter
	for _, p := range pairs {
		f, err := parseLengthPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseLengthPair(pair string) (*filter.LengthFilter, error) {
	parts := strings.SplitN(pair, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("length flag: expected field:min:max, got %q", pair)
	}
	field := parts[0]
	if field == "" {
		return nil, fmt.Errorf("length flag: field name must not be empty in %q", pair)
	}
	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("length flag: invalid min in %q: %w", pair, err)
	}
	max, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("length flag: invalid max in %q: %w", pair, err)
	}
	return filter.NewLengthFilter(field, min, max)
}
