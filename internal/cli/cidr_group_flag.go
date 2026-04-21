package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseCIDRGroupPairs parses repeated --cidr-group flags of the form
// "field=cidr1,cidr2,..." and returns a slice of CIDRGroupFilters.
func parseCIDRGroupPairs(pairs []string) ([]*filter.CIDRGroupFilter, error) {
	filters := make([]*filter.CIDRGroupFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseCIDRGroupPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// parseCIDRGroupPair parses a single "field=cidr1,cidr2" token.
func parseCIDRGroupPair(pair string) (*filter.CIDRGroupFilter, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return nil, fmt.Errorf("cidr-group: expected field=cidr1,cidr2 format, got %q", pair)
	}
	field := pair[:idx]
	rest := pair[idx+1:]

	if field == "" {
		return nil, fmt.Errorf("cidr-group: field name must not be empty in %q", pair)
	}
	if rest == "" {
		return nil, fmt.Errorf("cidr-group: CIDR list must not be empty in %q", pair)
	}

	cidrs := strings.Split(rest, ",")
	return filter.NewCIDRGroupFilter(field, cidrs)
}
