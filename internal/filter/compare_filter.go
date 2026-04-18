package filter

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// CompareOp represents a comparison operator.
type CompareOp string

const (
	OpEq  CompareOp = "eq"
	OpNeq CompareOp = "neq"
	OpGt  CompareOp = "gt"
	OpGte CompareOp = "gte"
	OpLt  CompareOp = "lt"
	OpLte CompareOp = "lte"
)

type compareFilter struct {
	field string
	op    CompareOp
	value float64
}

// NewCompareFilter returns a filter that compares a numeric field using op against value.
func NewCompareFilter(field string, op CompareOp, value float64) (*compareFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("compare filter: field must not be empty")
	}
	switch op {
	case OpEq, OpNeq, OpGt, OpGte, OpLt, OpLte:
	default:
		return nil, fmt.Errorf("compare filter: unknown operator %q", op)
	}
	return &compareFilter{field: field, op: op, value: value}, nil
}

func (f *compareFilter) Match(e parser.LogEntry) bool {
	raw, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	var num float64
	switch v := raw.(type) {
	case float64:
		num = v
	case string:
		var err error
		num, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
	default:
		return false
	}
	switch f.op {
	case OpEq:
		return num == f.value
	case OpNeq:
		return num != f.value
	case OpGt:
		return num > f.value
	case OpGte:
		return num >= f.value
	case OpLt:
		return num < f.value
	case OpLte:
		return num <= f.value
	}
	return false
}
