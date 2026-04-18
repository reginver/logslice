package filter

import (
	"fmt"
	"time"
)

// TimeFilter holds an optional start and end time for filtering log entries.
type TimeFilter struct {
	Start *time.Time
	End   *time.Time
}

// NewTimeFilter creates a TimeFilter from optional start and end time strings.
// Accepted format: RFC3339 (e.g. 2024-01-15T10:00:00Z)
func NewTimeFilter(start, end string) (*TimeFilter, error) {
	tf := &TimeFilter{}

	if start != "" {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return nil, fmt.Errorf("invalid start time %q: %w", start, err)
		}
		tf.Start = &t
	}

	if end != "" {
		t, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return nil, fmt.Errorf("invalid end time %q: %w", end, err)
		}
		tf.End = &t
	}

	if tf.Start != nil && tf.End != nil && tf.End.Before(*tf.Start) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	return tf, nil
}

// Match returns true if the given timestamp falls within the filter range.
func (tf *TimeFilter) Match(t time.Time) bool {
	if tf.Start != nil && t.Before(*tf.Start) {
		return false
	}
	if tf.End != nil && t.After(*tf.End) {
		return false
	}
	return true
}
