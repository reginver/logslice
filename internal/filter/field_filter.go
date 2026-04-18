package filter

import (
	"fmt"
	"regexp"
)

// FieldFilter matches log entries by a key=pattern pair.
type FieldFilter struct {
	Key     string
	Pattern *regexp.Regexp
}

// NewFieldFilter creates a FieldFilter for the given key and regex pattern.
func NewFieldFilter(key, pattern string) (*FieldFilter, error) {
	if key == "" {
		return nil, fmt.Errorf("field key must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	return &FieldFilter{Key: key, Pattern: re}, nil
}

// Match returns true if the entry map contains the key and its string value
// matches the compiled pattern.
func (ff *FieldFilter) Match(entry map[string]interface{}) bool {
	val, ok := entry[ff.Key]
	if !ok {
		return false
	}
	s := fmt.Sprintf("%v", val)
	return ff.Pattern.MatchString(s)
}

// FieldFilters is a slice of FieldFilter that all must match (AND semantics).
type FieldFilters []*FieldFilter

// Match returns true only if all filters match the entry.
func (ffs FieldFilters) Match(entry map[string]interface{}) bool {
	for _, f := range ffs {
		if !f.Match(entry) {
			return false
		}
	}
	return true
}
