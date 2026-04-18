package cli

import (
	"fmt"

	"github.com/user/logslice/internal/filter"
)

// parseCountFilter parses the --limit flag value and returns a CountFilter.
// Returns nil, nil when limit is zero (disabled).
func parseCountFilter(limit int) (*filter.CountFilter, error) {
	if limit == 0 {
		return nil, nil
	}
	f, err := filter.NewCountFilter(limit)
	if err != nil {
		return nil, fmt.Errorf("--limit: %w", err)
	}
	return f, nil
}
