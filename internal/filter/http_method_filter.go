package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// HTTPMethodFilter matches log entries whose HTTP method field is one of the
// allowed methods (e.g. GET, POST, PUT, DELETE).
type HTTPMethodFilter struct {
	field   string
	methods map[string]struct{}
}

// NewHTTPMethodFilter creates a filter that passes entries where the named
// field contains one of the given HTTP methods. Method names are normalised to
// upper-case before comparison.
func NewHTTPMethodFilter(field string, methods []string) (*HTTPMethodFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("http_method_filter: field must not be empty")
	}
	if len(methods) == 0 {
		return nil, fmt.Errorf("http_method_filter: at least one method is required")
	}
	allowed := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		m = strings.TrimSpace(m)
		if m == "" {
			return nil, fmt.Errorf("http_method_filter: method must not be empty")
		}
		allowed[strings.ToUpper(m)] = struct{}{}
	}
	return &HTTPMethodFilter{field: field, methods: allowed}, nil
}

// Match returns true when the entry's field value (case-insensitive) is one of
// the configured HTTP methods.
func (f *HTTPMethodFilter) Match(e *parser.Entry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	_, found := f.methods[strings.ToUpper(strings.TrimSpace(s))]
	return found
}
