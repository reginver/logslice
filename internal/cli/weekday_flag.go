package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseWeekdayPairs parses a slice of "field=day1,day2" strings and returns
// the corresponding WeekdayFilters. Each pair maps a field name to a
// comma-separated list of weekday names.
func parseWeekdayPairs(pairs []string) ([]*filter.WeekdayFilter, error) {
	filters := make([]*filter.WeekdayFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseWeekdayPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// parseWeekdayPair parses a single "field=day1,day2" token.
func parseWeekdayPair(pair string) (*filter.WeekdayFilter, error) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("weekday flag: expected field=day1,day2 format, got %q", pair)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("weekday flag: field name must not be empty in %q", pair)
	}
	rawDays := strings.TrimSpace(parts[1])
	if rawDays == "" {
		return nil, fmt.Errorf("weekday flag: weekday list must not be empty in %q", pair)
	}
	days := strings.Split(rawDays, ",")
	for i, d := range days {
		days[i] = strings.TrimSpace(d)
	}
	return filter.NewWeekdayFilter(field, days)
}
