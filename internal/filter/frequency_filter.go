package filter

import (
	"fmt"
	"sync"
)

// FrequencyFilter passes entries whose field value appears at least MinCount
// and at most MaxCount times across the entire stream.
type FrequencyFilter struct {
	field    string
	minCount int
	maxCount int
	mu       sync.Mutex
	counts   map[string]int
	buffer   []entry
}

type entry struct {
	e     map[string]any
	value string
}

// NewFrequencyFilter creates a FrequencyFilter for the given field.
// minCount and maxCount are inclusive bounds; set maxCount to -1 for unbounded.
func NewFrequencyFilter(field string, minCount, maxCount int) (*FrequencyFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("frequency filter: field must not be empty")
	}
	if minCount < 1 {
		return nil, fmt.Errorf("frequency filter: minCount must be >= 1")
	}
	if maxCount != -1 && maxCount < minCount {
		return nil, fmt.Errorf("frequency filter: maxCount (%d) must be >= minCount (%d)", maxCount, minCount)
	}
	return &FrequencyFilter{
		field:    field,
		minCount: minCount,
		maxCount: maxCount,
		counts:   make(map[string]int),
	}, nil
}

// Tally records the field value from the entry for later frequency evaluation.
func (f *FrequencyFilter) Tally(e map[string]any) {
	v, ok := e[f.field]
	if !ok {
		return
	}
	key := fmt.Sprintf("%v", v)
	f.mu.Lock()
	f.counts[key]++
	f.mu.Unlock()
}

// Match returns true if the entry's field value frequency falls within bounds.
func (f *FrequencyFilter) Match(e map[string]any) bool {
	v, ok := e[f.field]
	if !ok {
		return false
	}
	key := fmt.Sprintf("%v", v)
	f.mu.Lock()
	count := f.counts[key]
	f.mu.Unlock()
	if count < f.minCount {
		return false
	}
	if f.maxCount != -1 && count > f.maxCount {
		return false
	}
	return true
}
