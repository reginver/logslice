package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/user/logslice/internal/filter"
)

// parseLatencyPairs parses repeated --latency flags of the form
// "field=minMs-maxMs", e.g. "latency_ms=100-500". Either bound may be
// omitted to leave it unbounded: "latency_ms=-500" or "latency_ms=100-".
func parseLatencyPairs(pairs []string) ([]*filter.LatencyFilter, error) {
	var filters []*filter.LatencyFilter
	for _, p := range pairs {
		f, err := parseLatencyPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseLatencyPair(pair string) (*filter.LatencyFilter, error) {
	eq := strings.IndexByte(pair, '=')
	if eq < 0 {
		return nil, fmt.Errorf("latency flag: expected field=minMs-maxMs, got %q", pair)
	}
	field := pair[:eq]
	if field == "" {
		return nil, fmt.Errorf("latency flag: field name must not be empty in %q", pair)
	}
	rest := pair[eq+1:]
	dash := strings.IndexByte(rest, '-')
	if dash < 0 {
		return nil, fmt.Errorf("latency flag: expected minMs-maxMs in %q", pair)
	}
	minStr := rest[:dash]
	maxStr := rest[dash+1:]
	if minStr == "" && maxStr == "" {
		return nil, fmt.Errorf("latency flag: at least one bound required in %q", pair)
	}
	var minD, maxD time.Duration
	if minStr != "" {
		v, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return nil, fmt.Errorf("latency flag: invalid min %q in %q", minStr, pair)
		}
		minD = time.Duration(v) * time.Millisecond
	}
	if maxStr != "" {
		v, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return nil, fmt.Errorf("latency flag: invalid max %q in %q", maxStr, pair)
		}
		maxD = time.Duration(v) * time.Millisecond
	}
	return filter.NewLatencyFilter(field, minD, maxD)
}
