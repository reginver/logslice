package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseTagPairs parses a slice of "field=tag1,tag2,..." strings into TagFilters.
func parseTagPairs(pairs []string) ([]filter.Filter, error) {
	filters := make([]filter.Filter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseTagPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseTagPair(pair string) (filter.Filter, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return nil, fmt.Errorf("tag filter: expected 'field=tag1,tag2' format, got %q", pair)
	}
	field := strings.TrimSpace(pair[:idx])
	if field == "" {
		return nil, fmt.Errorf("tag filter: field name must not be empty in %q", pair)
	}
	raw := strings.TrimSpace(pair[idx+1:])
	if raw == "" {
		return nil, fmt.Errorf("tag filter: tag list must not be empty in %q", pair)
	}
	tags := strings.Split(raw, ",")
	return filter.NewTagFilter(field, tags)
}
