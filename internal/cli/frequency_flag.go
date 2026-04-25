package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/example/logslice/internal/filter"
)

// parseFrequencyPairs parses repeated --freq flags of the form "field:min:max".
// max may be omitted or set to "-1" for unbounded.
func parseFrequencyPairs(pairs []string) ([]*filter.FrequencyFilter, error) {
	filters := make([]*filter.FrequencyFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseFrequencyPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseFrequencyPair(s string) (*filter.FrequencyFilter, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("frequency flag: expected field:min[:max], got %q", s)
	}
	field := parts[0]
	if field == "" {
		return nil, fmt.Errorf("frequency flag: field must not be empty in %q", s)
	}
	minCount, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("frequency flag: invalid min count %q: %w", parts[1], err)
	}
	maxCount := -1
	if len(parts) == 3 && parts[2] != "" {
		maxCount, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("frequency flag: invalid max count %q: %w", parts[2], err)
		}
	}
	return filter.NewFrequencyFilter(field, minCount, maxCount)
}
