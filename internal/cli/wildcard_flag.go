package cli

import (
	"fmt"
	"strings"
)

type wildcardPair struct {
	field   string
	pattern string
}

// parseWildcardPairs parses a slice of "field=pattern" strings.
func parseWildcardPairs(args []string) ([]wildcardPair, error) {
	pairs := make([]wildcardPair, 0, len(args))
	for _, a := range args {
		p, err := parseWildcardPair(a)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func parseWildcardPair(s string) (wildcardPair, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return wildcardPair{}, fmt.Errorf("wildcard: expected field=pattern, got %q", s)
	}
	field, pattern := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if field == "" {
		return wildcardPair{}, fmt.Errorf("wildcard: field must not be empty in %q", s)
	}
	if pattern == "" {
		return wildcardPair{}, fmt.Errorf("wildcard: pattern must not be empty in %q", s)
	}
	return wildcardPair{field: field, pattern: pattern}, nil
}
