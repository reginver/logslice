package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// DateFilter matches entries whose timestamp falls on a specific date (YYYY-MM-DD).
type DateFilter struct {
	date  time.Time
	field string
}

// NewDateFilter returns a filter matching entries on the given date string (YYYY-MM-DD).
// If field is empty, it uses the entry's parsed Timestamp.
func NewDateFilter(date, field string) (*DateFilter, error) {
	if date == "" {
		return nil, fmt.Errorf("date filter: date must not be empty")
	}
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("date filter: invalid date %q: %w", date, err)
	}
	return &DateFilter{date: d, field: field}, nil
}

func (f *DateFilter) Match(e parser.Entry) bool {
	var t *time.Time
	if f.field != "" {
		v, ok := e.Fields[f.field]
		if !ok {
			return false
		}
		if s, ok := v.(string); ok {
			parsed, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return false
			}
			t = &parsed
		}
	} else {
		t = e.Timestamp
	}
	if t == nil {
		return false
	}
	y1, m1, d1 := t.Date()
	y2, m2, d2 := f.date.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
