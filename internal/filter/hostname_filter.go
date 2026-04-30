package filter

import (
	"fmt"
	"strings"
)

// HostnameFilter matches log entries whose hostname field value matches
// one of a set of allowed glob-style patterns (e.g. "web-*", "db-01").
type HostnameFilter struct {
	field    string
	patterns []string
}

// NewHostnameFilter returns a filter that passes entries where the value of
// field matches at least one of the provided patterns. Patterns support a
// single '*' wildcard that matches any sequence of characters.
func NewHostnameFilter(field string, patterns []string) (*HostnameFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("hostname filter: field must not be empty")
	}
	if len(patterns) == 0 {
		return nil, fmt.Errorf("hostname filter: at least one pattern required")
	}
	for _, p := range patterns {
		if p == "" {
			return nil, fmt.Errorf("hostname filter: pattern must not be empty")
		}
	}
	return &HostnameFilter{field: field, patterns: patterns}, nil
}

// Match returns true if the entry's field value matches any of the patterns.
func (f *HostnameFilter) Match(entry map[string]interface{}) bool {
	v, ok := entry[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	for _, p := range f.patterns {
		if hostnameGlobMatch(p, s) {
			return true
		}
	}
	return false
}

// hostnameGlobMatch matches s against a simple glob pattern where '*'
// matches any sequence of characters. Only a single '*' per pattern segment
// is supported.
func hostnameGlobMatch(pattern, s string) bool {
	parts := strings.SplitN(pattern, "*", 2)
	if len(parts) == 1 {
		return pattern == s
	}
	prefix, suffix := parts[0], parts[1]
	if !strings.HasPrefix(s, prefix) {
		return false
	}
	return strings.HasSuffix(s[len(prefix):], suffix)
}
