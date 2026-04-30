package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseChecksumPairs parses a slice of "field=algo:prefix" tokens into
// ChecksumFilter instances. Each token must contain exactly one '=' separating
// the field name from "algo:prefix".
func parseChecksumPairs(pairs []string) ([]*filter.ChecksumFilter, error) {
	var filters []*filter.ChecksumFilter
	for _, p := range pairs {
		f, err := parseChecksumPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// parseChecksumPair parses a single "field=algo:prefix" token.
func parseChecksumPair(s string) (*filter.ChecksumFilter, error) {
	eqIdx := strings.IndexByte(s, '=')
	if eqIdx < 0 {
		return nil, fmt.Errorf("checksum flag: expected field=algo:prefix, got %q", s)
	}
	field := s[:eqIdx]
	rest := s[eqIdx+1:]
	if field == "" {
		return nil, fmt.Errorf("checksum flag: field must not be empty in %q", s)
	}
	colIdx := strings.IndexByte(rest, ':')
	if colIdx < 0 {
		return nil, fmt.Errorf("checksum flag: expected algo:prefix after '=', got %q", rest)
	}
	algo := rest[:colIdx]
	prefix := rest[colIdx+1:]
	if algo == "" {
		return nil, fmt.Errorf("checksum flag: algo must not be empty in %q", s)
	}
	if prefix == "" {
		return nil, fmt.Errorf("checksum flag: prefix must not be empty in %q", s)
	}
	return filter.NewChecksumFilter(field, algo, prefix)
}
