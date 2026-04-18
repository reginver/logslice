package cli

import (
	"fmt"
	"strings"
)

// containsPair holds a parsed field=value1,value2 specification.
type containsPair struct {
	Field  string
	Values []string
}

// parseContainsPairs parses a slice of "field=val1,val2" strings into containsPair
// structs suitable for building ContainsFilters.
func parseContainsPairs(args []string) ([]containsPair, error) {
	pairs := make([]containsPair, 0, len(args))
	for _, arg := range args {
		p, err := parseContainsPair(arg)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func parseContainsPair(s string) (containsPair, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return containsPair{}, fmt.Errorf("contains: expected field=value[,value…], got %q", s)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return containsPair{}, fmt.Errorf("contains: field name must not be empty in %q", s)
	}
	raw := strings.Split(parts[1], ",")
	values := make([]string, 0, len(raw))
	for _, v := range raw {
		v = strings.TrimSpace(v)
		if v == "" {
			return containsPair{}, fmt.Errorf("contains: value must not be empty in %q", s)
		}
		values = append(values, v)
	}
	if len(values) == 0 {
		return containsPair{}, fmt.Errorf("contains: at least one value required in %q", s)
	}
	return containsPair{Field: field, Values: values}, nil
}
