package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// TenantFilter matches log entries whose tenant field value is one of the
// allowed tenant IDs. Matching is case-insensitive.
type TenantFilter struct {
	field   string
	allowed map[string]struct{}
}

// NewTenantFilter constructs a TenantFilter for the given field and tenant list.
// field must be non-empty and tenants must contain at least one non-empty value.
func NewTenantFilter(field string, tenants []string) (*TenantFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("tenant filter: field must not be empty")
	}
	if len(tenants) == 0 {
		return nil, fmt.Errorf("tenant filter: at least one tenant ID is required")
	}
	allowed := make(map[string]struct{}, len(tenants))
	for _, t := range tenants {
		if t == "" {
			return nil, fmt.Errorf("tenant filter: tenant ID must not be empty")
		}
		allowed[strings.ToLower(t)] = struct{}{}
	}
	return &TenantFilter{field: field, allowed: allowed}, nil
}

// Match returns true when the entry's field value (case-insensitive) is one of
// the configured tenant IDs.
func (f *TenantFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok || v == nil {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, found := f.allowed[strings.ToLower(s)]
	return found
}
