package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// parseHostnamePairs parses a slice of "field=pattern1,pattern2" strings into
// HostnameFilter instances. Multiple patterns for the same field are combined
// into a single filter.
func parseHostnamePairs(pairs []string) ([]*filter.HostnameFilter, error) {
	var filters []*filter.HostnameFilter
	for _, p := range pairs {
		f, err := parseHostnamePair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseHostnamePair(pair string) (*filter.HostnameFilter, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return nil, fmt.Errorf("hostname flag: expected field=pattern[,pattern...], got %q", pair)
	}
	field := pair[:idx]
	value := pair[idx+1:]
	if field == "" {
		return nil, fmt.Errorf("hostname flag: field must not be empty in %q", pair)
	}
	if value == "" {
		return nil, fmt.Errorf("hostname flag: pattern list must not be empty in %q", pair)
	}
	patterns := strings.Split(value, ",")
	return filter.NewHostnameFilter(field, patterns)
}
