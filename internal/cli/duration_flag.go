package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

// parseDurationPairs parses a slice of "field:min:max" strings into
// DurationFilter instances. Either min or max may be omitted (empty string)
// to leave that bound unset, but not both.
//
// Example: "latency:100ms:2s", "response_time::5s"
func parseDurationPairs(pairs []string) ([]*filter.DurationFilter, error) {
	var filters []*filter.DurationFilter
	for _, p := range pairs {
		f, err := parseDurationPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseDurationPair(pair string) (*filter.DurationFilter, error) {
	parts := strings.SplitN(pair, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("duration filter: expected field:min:max, got %q", pair)
	}
	field, rawMin, rawMax := parts[0], parts[1], parts[2]
	if field == "" {
		return nil, fmt.Errorf("duration filter: field must not be empty in %q", pair)
	}

	var minD, maxD time.Duration
	var err error
	if rawMin != "" {
		minD, err = time.ParseDuration(rawMin)
		if err != nil {
			return nil, fmt.Errorf("duration filter: invalid min %q: %w", rawMin, err)
		}
	}
	if rawMax != "" {
		maxD, err = time.ParseDuration(rawMax)
		if err != nil {
			return nil, fmt.Errorf("duration filter: invalid max %q: %w", rawMax, err)
		}
	}
	return filter.NewDurationFilter(field, minD, maxD)
}
