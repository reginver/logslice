package filter

import (
	"errors"

	"github.com/yourorg/logslice/internal/parser"
)

// NullFilter matches entries where a field is either absent or has a null/empty value.
type NullFilter struct {
	field  string
	invert bool
}

// NewNullFilter returns a filter that matches entries where field is null/missing.
// If invert is true, it matches entries where the field is present and non-empty.
func NewNullFilter(field string, invert bool) (*NullFilter, error) {
	if field == "" {
		return nil, errors.New("null_filter: field name must not be empty")
	}
	return &NullFilter{field: field, invert: invert}, nil
}

func (f *NullFilter) Match(e parser.LogEntry) bool {
	val, ok := e.Fields[f.field]
	isNull := !ok || val == nil || val == ""
	if f.invert {
		return !isNull
	}
	return isNull
}
