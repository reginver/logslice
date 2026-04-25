package filter

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// TimeOfDayFilter matches log entries whose timestamp falls within a
// wall-clock time-of-day window (e.g. 09:00–17:00), regardless of date.
type TimeOfDayFilter struct {
	start time.Duration // offset from midnight
	end   time.Duration
}

// NewTimeOfDayFilter creates a filter that passes entries whose time-of-day
// is in [start, end]. Both values must be in "HH:MM" or "HH:MM:SS" format.
func NewTimeOfDayFilter(start, end string) (*TimeOfDayFilter, error) {
	if start == "" || end == "" {
		return nil, fmt.Errorf("time_of_day_filter: start and end are required")
	}
	s, err := parseTOD(start)
	if err != nil {
		return nil, fmt.Errorf("time_of_day_filter: invalid start %q: %w", start, err)
	}
	e, err := parseTOD(end)
	if err != nil {
		return nil, fmt.Errorf("time_of_day_filter: invalid end %q: %w", end, err)
	}
	if s > e {
		return nil, fmt.Errorf("time_of_day_filter: start %q is after end %q", start, end)
	}
	return &TimeOfDayFilter{start: s, end: e}, nil
}

// Match returns true when the entry has a timestamp whose time-of-day is
// within [start, end].
func (f *TimeOfDayFilter) Match(e *parser.Entry) bool {
	if e.Timestamp == nil {
		return false
	}
	t := *e.Timestamp
	tod := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second
	return tod >= f.start && tod <= f.end
}

// parseTOD converts "HH:MM" or "HH:MM:SS" into a duration from midnight.
func parseTOD(s string) (time.Duration, error) {
	for _, layout := range []string{"15:04:05", "15:04"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return time.Duration(t.Hour())*time.Hour +
				time.Duration(t.Minute())*time.Minute +
				time.Duration(t.Second())*time.Second, nil
		}
	}
	return 0, fmt.Errorf("expected HH:MM or HH:MM:SS")
}
