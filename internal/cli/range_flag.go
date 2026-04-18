package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// RangePair holds a parsed field + numeric range specification.
type RangePair struct {
	Field string
	Min   float64
	Max   float64
}

// parseRangePairs parses slice of "field:min:max" strings into RangePair values.
func parseRangePairs(specs []string) ([]RangePair, error) {
	pairs := make([]RangePair, 0, len(specs))
	for _, s := range specs {
		p, err := parseRangePair(s)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func parseRangePair(s string) (RangePair, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return RangePair{}, fmt.Errorf("range: invalid format %q, expected field:min:max", s)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return RangePair{}, fmt.Errorf("range: field name is empty in %q", s)
	}
	min, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return RangePair{}, fmt.Errorf("range: invalid min in %q: %w", s, err)
	}
	max, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return RangePair{}, fmt.Errorf("range: invalid max in %q: %w", s, err)
	}
	return RangePair{Field: field, Min: min, Max: max}, nil
}
