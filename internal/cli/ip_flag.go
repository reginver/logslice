package cli

import (
	"fmt"
	"strings"

	"github.com/your/logslice/internal/filter"
)

// parseIPPairs parses a slice of "field=cidr" strings into IPFilters.
func parseIPPairs(pairs []string) ([]*filter.IPFilter, error) {
	filters := make([]*filter.IPFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseIPPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseIPPair(pair string) (*filter.IPFilter, error) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("ip flag: expected field=cidr, got %q", pair)
	}
	field := strings.TrimSpace(parts[0])
	cidr := strings.TrimSpace(parts[1])
	if field == "" {
		return nil, fmt.Errorf("ip flag: field name must not be empty in %q", pair)
	}
	if cidr == "" {
		return nil, fmt.Errorf("ip flag: cidr must not be empty in %q", pair)
	}
	return filter.NewIPFilter(field, cidr)
}
