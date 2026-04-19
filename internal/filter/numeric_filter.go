package filter

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// NumericFilter matches entries where a numeric field is odd or even,
// or more usefully — matches entries by divisibility (modulo check).
// More practically: this filter checks if a numeric field's value
// satisfies value % divisor == remainder.
type NumericFilter struct {
	field     string
	divisor   float64
	remainder float64
}

// NewNumericFilter creates a filter that matches entries where
// field value mod divisor == remainder.
func NewNumericFilter(field string, divisor, remainder float64) (*NumericFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("numeric filter: field must not be empty")
	}
	if divisor == 0 {
		return nil, fmt.Errorf("numeric filter: divisor must not be zero")
	}
	if remainder < 0 || remainder >= divisor {
		return nil, fmt.Errorf("numeric filter: remainder must be in [0, divisor)")
	}
	return &NumericFilter{field: field, divisor: divisor, remainder: remainder}, nil
}

func (f *NumericFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	var num float64
	switch val := v.(type) {
	case float64:
		num = val
	case int:
		num = float64(val)
	case string:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return false
		}
		num = parsed
	default:
		return false
	}
	return math_mod(num, f.divisor) == f.remainder
}

func math_mod(a, b float64) float64 {
	result := a - float64(int(a/b))*b
	if result < 0 {
		result += b
	}
	return result
}
