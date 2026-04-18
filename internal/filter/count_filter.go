package filter

import (
	"fmt"
	"sync/atomic"

	"github.com/user/logslice/internal/parser"
)

// CountFilter passes only the first N matched entries.
type CountFilter struct {
	max   int64
	count atomic.Int64
}

// NewCountFilter returns a filter that passes at most max entries.
// max must be greater than zero.
func NewCountFilter(max int) (*CountFilter, error) {
	if max <= 0 {
		return nil, fmt.Errorf("count filter: max must be greater than zero, got %d", max)
	}
	return &CountFilter{max: int64(max)}, nil
}

func (f *CountFilter) Match(e parser.LogEntry) bool {
	next := f.count.Add(1)
	return next <= f.max
}
