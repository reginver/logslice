package filter

import (
	"fmt"
	"strings"
)

// UserAgentFilter matches log entries whose user-agent field contains one of
// the specified browser/client tokens (case-insensitive substring match).
type UserAgentFilter struct {
	field  string
	tokens []string
}

// NewUserAgentFilter returns a filter that passes entries where the value of
// field contains at least one of the provided tokens (case-insensitive).
// field must not be empty and at least one non-empty token must be supplied.
func NewUserAgentFilter(field string, tokens []string) (*UserAgentFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("user_agent_filter: field must not be empty")
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("user_agent_filter: at least one token required")
	}
	norm := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if t == "" {
			return nil, fmt.Errorf("user_agent_filter: token must not be empty")
		}
		norm = append(norm, strings.ToLower(t))
	}
	return &UserAgentFilter{field: field, tokens: norm}, nil
}

// Match returns true when the entry's field value (lowercased) contains at
// least one of the filter tokens.
func (f *UserAgentFilter) Match(entry map[string]interface{}) bool {
	v, ok := entry[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	lower := strings.ToLower(s)
	for _, tok := range f.tokens {
		if strings.Contains(lower, tok) {
			return true
		}
	}
	return false
}
