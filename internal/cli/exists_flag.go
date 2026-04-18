package cli

import (
	"fmt"
	"strings"
)

// existsField holds a parsed field name for the exists filter.
type existsField struct {
	Field string
}

// parseExistsFields parses a list of field names for the --exists flag.
// Each value should be a non-empty field name.
func parseExistsFields(values []string) ([]existsField, error) {
	result := make([]existsField, 0, len(values))
	for _, v := range values {
		f, err := parseExistsField(v)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}
	return result, nil
}

// parseExistsField validates and returns a single existsField.
func parseExistsField(value string) (existsField, error) {
	field := strings.TrimSpace(value)
	if field == "" {
		return existsField{}, fmt.Errorf("exists: field name must not be empty")
	}
	return existsField{Field: field}, nil
}
