package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseErrorPairs parses a slice of "field=code1,code2,..." strings into
// ErrorFilter instances. Each pair must contain exactly one '=' separator.
func parseErrorPairs(pairs []string) ([]*filter.ErrorFilter, error) {
	filters := make([]*filter.ErrorFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseErrorPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

func parseErrorPair(pair string) (*filter.ErrorFilter, error) {
	idx := strings.IndexByte(pair, '=')
	if idx < 0 {
		return nil, fmt.Errorf("error filter: expected 'field=code1,code2' got %q", pair)
	}
	field := strings.TrimSpace(pair[:idx])
	raw := strings.TrimSpace(pair[idx+1:])
	if field == "" {
		return nil, fmt.Errorf("error filter: field name must not be empty in %q", pair)
	}
	if raw == "" {
		return nil, fmt.Errorf("error filter: code list must not be empty in %q", pair)
	}
	codes := strings.Split(raw, ",")
	return filter.NewErrorFilter(field, codes)
}
