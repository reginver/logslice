package cli

import (
	"fmt"
	"strings"
)

type prefixPair struct {
	Field  string
	Prefix string
}

func parsePrefixPairs(args []string) ([]prefixPair, error) {
	var pairs []prefixPair
	for _, arg := range args {
		p, err := parsePrefixPair(arg)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func parsePrefixPair(s string) (prefixPair, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return prefixPair{}, fmt.Errorf("invalid prefix filter %q: expected field=prefix", s)
	}
	field := strings.TrimSpace(parts[0])
	prefix := strings.TrimSpace(parts[1])
	if field == "" {
		return prefixPair{}, fmt.Errorf("invalid prefix filter %q: field must not be empty", s)
	}
	if prefix == "" {
		return prefixPair{}, fmt.Errorf("invalid prefix filter %q: prefix must not be empty", s)
	}
	return prefixPair{Field: field, Prefix: prefix}, nil
}
