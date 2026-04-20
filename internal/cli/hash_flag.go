package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseHashPairs parses a slice of "field=prefix" strings into HashFilter instances.
// Each element must be in the form "field=hexprefix".
func parseHashPairs(pairs []string) ([]*filter.HashFilter, error) {
	filters := make([]*filter.HashFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseHashPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// parseHashPair parses a single "field=hexprefix" string.
func parseHashPair(pair string) (*filter.HashFilter, error) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("hash filter: invalid format %q, expected field=prefix", pair)
	}
	field := strings.TrimSpace(parts[0])
	prefix := strings.TrimSpace(parts[1])
	if field == "" {
		return nil, fmt.Errorf("hash filter: field name must not be empty in %q", pair)
	}
	if prefix == "" {
		return nil, fmt.Errorf("hash filter: prefix must not be empty in %q", pair)
	}
	return filter.NewHashFilter(field, prefix)
}
