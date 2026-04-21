package filter

import (
	"fmt"
	"strings"
)

// EnvFilter matches log entries where a string field value matches one of the
// allowed environment names (e.g. "production", "staging", "development").
// Comparison is case-insensitive.
type EnvFilter struct {
	field   string
	allowed map[string]struct{}
}

// NewEnvFilter creates an EnvFilter for the given field and list of environment
// names. At least one environment name must be provided, and neither the field
// nor any environment name may be empty.
func NewEnvFilter(field string, envs []string) (*EnvFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("env filter: field must not be empty")
	}
	if len(envs) == 0 {
		return nil, fmt.Errorf("env filter: at least one environment name required")
	}
	allowed := make(map[string]struct{}, len(envs))
	for _, e := range envs {
		if e == "" {
			return nil, fmt.Errorf("env filter: environment name must not be empty")
		}
		allowed[strings.ToLower(e)] = struct{}{}
	}
	return &EnvFilter{field: field, allowed: allowed}, nil
}

// Match returns true when the entry's field value (case-insensitive) is one of
// the allowed environment names.
func (f *EnvFilter) Match(entry map[string]interface{}) bool {
	v, ok := entry[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, found := f.allowed[strings.ToLower(s)]
	return found
}
