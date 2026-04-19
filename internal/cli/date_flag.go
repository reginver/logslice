package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// parseDatePairs parses --date flags of the form "YYYY-MM-DD" or "field=YYYY-MM-DD".
func parseDatePairs(values []string) ([]*filter.DateFilter, error) {
	filters := make([]*filter.DateFilter, 0, len(values))
	for _, v := range values {
		f, err := parseDatePair(v)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseDatePair(s string) (*filter.DateFilter, error) {
	if s == "" {
		return nil, fmt.Errorf("date flag: value must not be empty")
	}
	if idx := strings.Index(s, "="); idx >= 0 {
		field := strings.TrimSpace(s[:idx])
		date := strings.TrimSpace(s[idx+1:])
		if field == "" {
			return nil, fmt.Errorf("date flag: field must not be empty in %q", s)
		}
		if date == "" {
			return nil, fmt.Errorf("date flag: date must not be empty in %q", s)
		}
		return filter.NewDateFilter(date, field)
	}
	// bare date — match against entry timestamp
	return filter.NewDateFilter(s, "")
}
