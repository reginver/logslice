package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// parseComparePairs parses a slice of "field:op:value" strings into CompareFilters.
func parseComparePairs(pairs []string) ([]*filter.compareFilter, error) {
	var filters []*filter.compareFilter
	for _, p := range pairs {
		f, err := parseComparePair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseComparePair(s string) (*filter.compareFilter, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("compare: expected field:op:value, got %q", s)
	}
	field, opStr, valStr := parts[0], parts[1], parts[2]
	if field == "" {
		return nil, fmt.Errorf("compare: field must not be empty in %q", s)
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, fmt.Errorf("compare: invalid value %q in %q", valStr, s)
	}
	return filter.NewCompareFilter(field, filter.CompareOp(opStr), val)
}
