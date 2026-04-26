package filter

import (
	"errors"
	"sync/atomic"

	"github.com/user/logslice/internal/parser"
)

// SampleFilter passes every Nth log entry, discarding the rest.
// N must be >= 1; N=1 passes all entries.
type SampleFilter struct {
	n       uint64
	counter atomic.Uint64
}

// NewSampleFilter creates a SampleFilter that passes one out of every n entries.
func NewSampleFilter(n uint64) (*SampleFilter, error) {
	if n == 0 {
		return nil, errors.New("sample: n must be >= 1")
	}
	return &SampleFilter{n: n}, nil
}

// Match returns true for every n-th call (1-indexed).
func (f *SampleFilter) Match(entry parser.LogEntry) bool {
	c := f.counter.Add(1)
	return c%f.n == 0
}
