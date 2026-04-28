package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// tagFilter matches log entries where a field contains all specified tags
// (comma-separated values in the field are treated as a tag set).
type tagFilter struct {
	field string
	tags  []string
}

// NewTagFilter returns a Filter that passes entries whose field value contains
// every tag in the provided list. Tags are matched case-insensitively and the
// field value is split on commas.
func NewTagFilter(field string, tags []string) (Filter, error) {
	if field == "" {
		return nil, fmt.Errorf("tag filter: field must not be empty")
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("tag filter: at least one tag is required")
	}
	norm := make([]string, len(tags))
	for i, t := range tags {
		t = strings.TrimSpace(t)
		if t == "" {
			return nil, fmt.Errorf("tag filter: tag must not be empty")
		}
		norm[i] = strings.ToLower(t)
	}
	return &tagFilter{field: field, tags: norm}, nil
}

func (f *tagFilter) Match(e parser.LogEntry) bool {
	raw, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := raw.(string)
	if !ok {
		return false
	}
	parts := strings.Split(s, ",")
	set := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		set[strings.ToLower(strings.TrimSpace(p))] = struct{}{}
	}
	for _, t := range f.tags {
		if _, found := set[t]; !found {
			return false
		}
	}
	return true
}
