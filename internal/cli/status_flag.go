package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseStatusPairs parses a slice of "field=range1,range2,..." strings and
// returns a slice of StatusFilters, one per pair.
//
// Example input: ["status=200-299,404", "http_code=500-503"]
func parseStatusPairs(pairs []string) ([]*filter.StatusFilter, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	out := make([]*filter.StatusFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseStatusPair(p)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func parseStatusPair(s string) (*filter.StatusFilter, error) {
	idx := strings.IndexByte(s, '=')
	if idx < 0 {
		return nil, fmt.Errorf("status flag: expected field=ranges, got %q", s)
	}
	field := strings.TrimSpace(s[:idx])
	if field == "" {
		return nil, fmt.Errorf("status flag: empty field in %q", s)
	}
	raw := strings.TrimSpace(s[idx+1:])
	if raw == "" {
		return nil, fmt.Errorf("status flag: empty range list in %q", s)
	}
	specs := strings.Split(raw, ",")
	for i, sp := range specs {
		specs[i] = strings.TrimSpace(sp)
	}
	return filter.NewStatusFilter(field, specs)
}
