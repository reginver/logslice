package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// KeywordFilter matches log entries whose target field contains at least one
// of the provided keywords (case-insensitive).
type KeywordFilter struct {
	field    string
	keywords []string
}

// NewKeywordFilter creates a KeywordFilter for the given field and keyword list.
// At least one keyword must be supplied and neither the field nor any keyword
// may be empty.
func NewKeywordFilter(field string, keywords []string) (*KeywordFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("keyword filter: field must not be empty")
	}
	if len(keywords) == 0 {
		return nil, fmt.Errorf("keyword filter: at least one keyword required")
	}
	norm := make([]string, 0, len(keywords))
	for _, k := range keywords {
		if k == "" {
			return nil, fmt.Errorf("keyword filter: keyword must not be empty")
		}
		norm = append(norm, strings.ToLower(k))
	}
	return &KeywordFilter{field: field, keywords: norm}, nil
}

// Match returns true when the entry's field value (as a string) contains at
// least one of the filter's keywords (case-insensitive).
func (f *KeywordFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s := strings.ToLower(fmt.Sprintf("%v", v))
	for _, k := range f.keywords {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}
