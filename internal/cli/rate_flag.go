package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/logslice/logslice/internal/filter"
)

// parseRatePairs parses repeated --rate flags of the form:
//
//	field=window,min,max
//
// where window is a Go duration string (e.g. "5s"), and min/max are
// floats per second (use 0 to leave a bound open).
func parseRatePairs(pairs []string) ([]*filter.RateFilter, error) {
	var filters []*filter.RateFilter
	for _, p := range pairs {
		f, err := parseRatePair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseRatePair(s string) (*filter.RateFilter, error) {
	eq := strings.IndexByte(s, '=')
	if eq < 0 {
		return nil, fmt.Errorf("rate flag: expected field=window,min,max, got %q", s)
	}
	field := strings.TrimSpace(s[:eq])
	if field == "" {
		return nil, fmt.Errorf("rate flag: field must not be empty in %q", s)
	}
	parts := strings.SplitN(strings.TrimSpace(s[eq+1:]), ",", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("rate flag: expected window,min,max after '=' in %q", s)
	}
	window, err := time.ParseDuration(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("rate flag: invalid window %q in %q: %w", parts[0], s, err)
	}
	minRate, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("rate flag: invalid min %q in %q: %w", parts[1], s, err)
	}
	maxRate, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return nil, fmt.Errorf("rate flag: invalid max %q in %q: %w", parts[2], s, err)
	}
	return filter.NewRateFilter(field, window, minRate, maxRate)
}
