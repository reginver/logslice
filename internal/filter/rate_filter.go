package filter

import (
	"fmt"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// RateFilter passes entries whose numeric field value falls within a
// [min, max] rate per second computed against the entry's timestamp.
// It accumulates a sliding window of values and computes rate on flush.
type RateFilter struct {
	field    string
	window   time.Duration
	minRate  float64
	maxRate  float64
	bucket   []ratePoint
}

type ratePoint struct {
	at    time.Time
	value float64
}

// NewRateFilter creates a RateFilter.
// field is the numeric log field, window is the aggregation window,
// minRate/maxRate are inclusive bounds (use 0 to leave a bound open).
func NewRateFilter(field string, window time.Duration, minRate, maxRate float64) (*RateFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("rate filter: field must not be empty")
	}
	if window <= 0 {
		return nil, fmt.Errorf("rate filter: window must be positive")
	}
	if minRate < 0 || maxRate < 0 {
		return nil, fmt.Errorf("rate filter: rates must be non-negative")
	}
	if maxRate > 0 && minRate > maxRate {
		return nil, fmt.Errorf("rate filter: minRate %.4f exceeds maxRate %.4f", minRate, maxRate)
	}
	return &RateFilter{
		field:   field,
		window:  window,
		minRate: minRate,
		maxRate: maxRate,
	}, nil
}

// Match returns true when the per-second rate of field values within the
// sliding window satisfies the configured bounds.
func (f *RateFilter) Match(e parser.LogEntry) bool {
	if e.Timestamp == nil {
		return false
	}
	v, ok := toFloat64(e.Fields[f.field])
	if !ok {
		return false
	}
	now := *e.Timestamp
	cutoff := now.Add(-f.window)

	// Append current point.
	f.bucket = append(f.bucket, ratePoint{at: now, value: v})

	// Evict points outside the window.
	filtered := f.bucket[:0]
	for _, p := range f.bucket {
		if !p.at.Before(cutoff) {
			filtered = append(filtered, p)
		}
	}
	f.bucket = filtered

	if len(f.bucket) < 2 {
		return false
	}

	// Sum values and compute elapsed seconds.
	var sum float64
	for _, p := range f.bucket {
		sum += p.value
	}
	elapsed := f.bucket[len(f.bucket)-1].at.Sub(f.bucket[0].at).Seconds()
	if elapsed <= 0 {
		return false
	}
	rate := sum / elapsed

	if f.minRate > 0 && rate < f.minRate {
		return false
	}
	if f.maxRate > 0 && rate > f.maxRate {
		return false
	}
	return true
}
