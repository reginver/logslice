package cli

import (
	"fmt"
	"strings"
)

// parseSuffixPairs parses a slice of "field=suffix" strings into field/suffix pairs.
func parseSuffixPairs(args []string) ([][2]string, error) {
	pairs := make([][2]string, 0, len(args))
	for _, a := range args {
		p, err := parseSuffixPair(a)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

// parseSuffixPair parses a single "field=suffix" string.
func parseSuffixPair(s string) ([2]string, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return [2]string{}, fmt.Errorf("suffix flag: invalid format %q, expected field=suffix", s)
	}
	field := strings.TrimSpace(parts[0])
	suffix := parts[1]
	if field == "" {
		return [2]string{}, fmt.Errorf("suffix flag: field must not be empty in %q", s)
	}
	if suffix == "" {
		return [2]string{}, fmt.Errorf("suffix flag: suffix must not be empty in %q", s)
	}
	return [2]string{field, suffix}, nil
}
