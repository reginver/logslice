package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/user/logslice/internal/filter"
)

// parseBurstPairs parses repeated --burst flags of the form:
//
//	field=<window>/<minCount>   e.g.  ip=1m/5
func parseBurstPairs(pairs []string) ([]*filter.BurstFilter, error) {
	var filters []*filter.BurstFilter
	for _, p := range pairs {
		f, err := parseBurstPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseBurstPair(pair string) (*filter.BurstFilter, error) {
	eqIdx := strings.IndexByte(pair, '=')
	if eqIdx < 0 {
		return nil, fmt.Errorf("burst: expected field=window/minCount, got %q", pair)
	}
	field := pair[:eqIdx]
	rest := pair[eqIdx+1:]
	if field == "" {
		return nil, fmt.Errorf("burst: field must not be empty")
	}

	slashIdx := strings.LastIndexByte(rest, '/')
	if slashIdx < 0 {
		return nil, fmt.Errorf("burst: expected window/minCount, got %q", rest)
	}
	windowStr := rest[:slashIdx]
	countStr := rest[slashIdx+1:]

	if windowStr == "" {
		return nil, fmt.Errorf("burst: window must not be empty")
	}
	if countStr == "" {
		return nil, fmt.Errorf("burst: minCount must not be empty")
	}

	window, err := time.ParseDuration(windowStr)
	if err != nil {
		return nil, fmt.Errorf("burst: invalid window %q: %w", windowStr, err)
	}

	minCount, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, fmt.Errorf("burst: invalid minCount %q: %w", countStr, err)
	}

	return filter.NewBurstFilter(field, window, minCount)
}
