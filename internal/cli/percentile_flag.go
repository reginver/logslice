package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// parsePercentilePairs parses repeated --percentile flags of the form
// "field=minPct-maxPct", e.g. "latency=10-90".
func parsePercentilePairs(pairs []string) ([]*filter.PercentileFilter, error) {
	var filters []*filter.PercentileFilter
	for _, p := range pairs {
		f, err := parsePercentilePair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parsePercentilePair(s string) (*filter.PercentileFilter, error) {
	eq := strings.IndexByte(s, '=')
	if eq < 0 {
		return nil, fmt.Errorf("percentile: expected field=min-max, got %q", s)
	}
	field := s[:eq]
	rest := s[eq+1:]
	if field == "" {
		return nil, fmt.Errorf("percentile: field must not be empty in %q", s)
	}
	dash := strings.IndexByte(rest, '-')
	if dash < 0 {
		return nil, fmt.Errorf("percentile: expected min-max range in %q", s)
	}
	minStr := rest[:dash]
	maxStr := rest[dash+1:]
	if minStr == "" || maxStr == "" {
		return nil, fmt.Errorf("percentile: min and max must not be empty in %q", s)
	}
	minPct, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		return nil, fmt.Errorf("percentile: invalid min %q: %w", minStr, err)
	}
	maxPct, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		return nil, fmt.Errorf("percentile: invalid max %q: %w", maxStr, err)
	}
	return filter.NewPercentileFilter(field, minPct, maxPct)
}
