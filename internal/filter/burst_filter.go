package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// BurstFilter passes entries where the same field value appears at least
// minCount times within the given window duration.
type BurstFilter struct {
	field    string
	window   time.Duration
	minCount int
	// buckets maps field value -> list of timestamps seen
	buckets map[string][]time.Time
}

// NewBurstFilter creates a BurstFilter that matches entries whose field value
// has been seen at least minCount times within window.
func NewBurstFilter(field string, window time.Duration, minCount int) (*BurstFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("burst filter: field must not be empty")
	}
	if window <= 0 {
		return nil, fmt.Errorf("burst filter: window must be positive")
	}
	if minCount < 2 {
		return nil, fmt.Errorf("burst filter: minCount must be >= 2")
	}
	return &BurstFilter{
		field:    field,
		window:   window,
		minCount: minCount,
		buckets:  make(map[string][]time.Time),
	}, nil
}

// Match returns true when the entry's field value has appeared at least
// minCount times (including this entry) within the sliding window.
func (f *BurstFilter) Match(e parser.Entry) bool {
	raw, ok := e.Fields[f.field]
	if !ok || raw == nil {
		return false
	}
	key := fmt.Sprintf("%v", raw)

	var now time.Time
	if e.Timestamp != nil {
		now = *e.Timestamp
	} else {
		now = time.Now()
	}

	cutoff := now.Add(-f.window)
	times := f.buckets[key]

	// evict entries outside the window
	filtered := times[:0]
	for _, t := range times {
		if !t.Before(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	f.buckets[key] = filtered

	return len(filtered) >= f.minCount
}
