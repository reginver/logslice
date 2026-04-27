package filter

import (
	"fmt"

	"github.com/yourorg/logslice/internal/parser"
)

// jsonTypeFilter matches log entries where a field's JSON type matches the expected type.
// Supported types: "string", "number", "bool", "null", "array", "object".
type jsonTypeFilter struct {
	field    string
	expected string
}

var validJSONTypes = map[string]bool{
	"string": true,
	"number": true,
	"bool":   true,
	"null":   true,
	"array":  true,
	"object": true,
}

// NewJSONTypeFilter returns a filter that matches entries where the named field
// holds a value of the given JSON type.
func NewJSONTypeFilter(field, typeName string) (*jsonTypeFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("json_type filter: field must not be empty")
	}
	if !validJSONTypes[typeName] {
		return nil, fmt.Errorf("json_type filter: unknown type %q; must be one of string, number, bool, null, array, object", typeName)
	}
	return &jsonTypeFilter{field: field, expected: typeName}, nil
}

func (f *jsonTypeFilter) Match(e parser.LogEntry) bool {
	val, ok := e.Fields[f.field]
	if !ok {
		return f.expected == "null"
	}
	return jsonTypeOf(val) == f.expected
}

// jsonTypeOf returns the JSON type name for a value decoded from a JSON object.
// It mirrors the type set supported by jsonTypeFilter.
func jsonTypeOf(val interface{}) string {
	switch val.(type) {
	case nil:
		return "null"
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "bool"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return ""
	}
}
