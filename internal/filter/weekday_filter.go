package filter

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// WeekdayFilter matches log entries whose timestamp falls on one of the
// specified weekdays (e.g. "monday", "tuesday", …).
type WeekdayFilter struct {
	field   string
	weekdays map[time.Weekday]struct{}
}

var weekdayNames = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

// NewWeekdayFilter returns a filter that passes entries whose named time field
// falls on any of the supplied weekday names. Names are case-insensitive.
func NewWeekdayFilter(field string, days []string) (*WeekdayFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("weekday filter: field name must not be empty")
	}
	if len(days) == 0 {
		return nil, fmt.Errorf("weekday filter: at least one weekday must be specified")
	}
	set := make(map[time.Weekday]struct{}, len(days))
	for _, d := range days {
		wd, ok := weekdayNames[strings.ToLower(strings.TrimSpace(d))]
		if !ok {
			return nil, fmt.Errorf("weekday filter: unknown weekday %q", d)
		}
		set[wd] = struct{}{}
	}
	return &WeekdayFilter{field: field, weekdays: set}, nil
}

// Match returns true when the entry's field value is a time whose weekday is
// in the allowed set.
func (f *WeekdayFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	var t time.Time
	switch val := v.(type) {
	case time.Time:
		t = val
	case string:
		var err error
		t, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return false
		}
	default:
		return false
	}
	_, ok = f.weekdays[t.Weekday()]
	return ok
}
