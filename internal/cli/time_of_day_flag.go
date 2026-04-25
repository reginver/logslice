package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseTimeOfDayPairs parses one or more "HH:MM-HH:MM" window specs from the
// --tod flag (repeatable) and returns a slice of TimeOfDayFilters.
// Each value must be in the form "START-END" where START and END are
// "HH:MM" or "HH:MM:SS".
func parseTimeOfDayPairs(values []string) ([]*filter.TimeOfDayFilter, error) {
	var filters []*filter.TimeOfDayFilter
	for _, v := range values {
		f, err := parseTimeOfDayPair(v)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// parseTimeOfDayPair parses a single "HH:MM-HH:MM" or "HH:MM:SS-HH:MM:SS"
// string into a TimeOfDayFilter.
func parseTimeOfDayPair(s string) (*filter.TimeOfDayFilter, error) {
	// Split on the last '-' that separates the two times. Times may contain
	// '-' only if someone passes negative offsets, which we don't support,
	// so we look for the separator between the two time tokens.
	// Format: "HH:MM-HH:MM" or "HH:MM:SS-HH:MM:SS".
	// We split on '-' but need to be careful: "09:00-17:00" has one '-'.
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("time_of_day: expected START-END, got %q", s)
	}
	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])
	if start == "" {
		return nil, fmt.Errorf("time_of_day: empty start in %q", s)
	}
	if end == "" {
		return nil, fmt.Errorf("time_of_day: empty end in %q", s)
	}
	return filter.NewTimeOfDayFilter(start, end)
}
