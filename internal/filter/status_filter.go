package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// StatusFilter matches log entries whose HTTP status code field falls within
// one or more ranges (e.g. "200-299", "404", "500-503").
type StatusFilter struct {
	field  string
	ranges [][2]int
}

// NewStatusFilter creates a StatusFilter for the given field and range specs.
// Each spec is either a single code ("404") or a dash-separated range ("200-299").
func NewStatusFilter(field string, specs []string) (*StatusFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("status filter: field must not be empty")
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("status filter: at least one status range is required")
	}
	ranges := make([][2]int, 0, len(specs))
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			return nil, fmt.Errorf("status filter: empty status spec")
		}
		parts := strings.SplitN(spec, "-", 2)
		lo, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("status filter: invalid code %q: %w", parts[0], err)
		}
		hi := lo
		if len(parts) == 2 {
			hi, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("status filter: invalid code %q: %w", parts[1], err)
			}
		}
		if lo > hi {
			return nil, fmt.Errorf("status filter: range %d-%d is invalid (lo > hi)", lo, hi)
		}
		ranges = append(ranges, [2]int{lo, hi})
	}
	return &StatusFilter{field: field, ranges: ranges}, nil
}

// Match returns true when the entry's field value is an integer that falls
// within at least one of the configured status ranges.
func (f *StatusFilter) Match(entry map[string]interface{}) bool {
	v, ok := entry[f.field]
	if !ok {
		return false
	}
	var code int
	switch val := v.(type) {
	case float64:
		code = int(val)
	case int:
		code = val
	case string:
		n, err := strconv.Atoi(val)
		if err != nil {
			return false
		}
		code = n
	default:
		return false
	}
	for _, r := range f.ranges {
		if code >= r[0] && code <= r[1] {
			return true
		}
	}
	return false
}
