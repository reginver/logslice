package cli

import (
	"fmt"
	"strings"
)

// SubstringSpec holds parsed parameters for a substring filter.
type SubstringSpec struct {
	Field     string
	Substring string
	CaseFold  bool
}

// parseSubstringPairs parses a slice of "field:substring" or "field:substring:i"
// tokens into SubstringSpec values. The optional ":i" suffix enables case-insensitive matching.
func parseSubstringPairs(pairs []string) ([]SubstringSpec, error) {
	specs := make([]SubstringSpec, 0, len(pairs))
	for _, p := range pairs {
		spec, err := parseSubstringPair(p)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func parseSubstringPair(pair string) (SubstringSpec, error) {
	parts := strings.SplitN(pair, ":", 3)
	if len(parts) < 2 {
		return SubstringSpec{}, fmt.Errorf("substring flag: expected field:substring, got %q", pair)
	}
	field := parts[0]
	substring := parts[1]
	if field == "" {
		return SubstringSpec{}, fmt.Errorf("substring flag: field must not be empty in %q", pair)
	}
	if substring == "" {
		return SubstringSpec{}, fmt.Errorf("substring flag: substring must not be empty in %q", pair)
	}
	caseFold := len(parts) == 3 && strings.ToLower(parts[2]) == "i"
	return SubstringSpec{Field: field, Substring: substring, CaseFold: caseFold}, nil
}
