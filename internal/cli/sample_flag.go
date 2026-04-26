package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/filter"
)

// parseSampleFilter parses the --sample flag value and returns a SampleFilter.
// The value must be a positive integer string, e.g. "3" to keep every 3rd entry.
func parseSampleFilter(raw string) (*filter.SampleFilter, error) {
	if raw == "" {
		return nil, errors.New("sample: value must not be empty")
	}
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("sample: invalid integer %q: %w", raw, err)
	}
	f, err := filter.NewSampleFilter(n)
	if err != nil {
		return nil, fmt.Errorf("sample: %w", err)
	}
	return f, nil
}
