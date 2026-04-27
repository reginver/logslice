package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseKeywordPairs parses a slice of "field=kw1,kw2,..." strings into
// KeywordFilter instances.
func parseKeywordPairs(pairs []string) ([]*filter.KeywordFilter, error) {
	filters := make([]*filter.KeywordFilter, 0, len(pairs))
	for _, p := range pairs {
		f, err := parseKeywordPair(p)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// parseKeywordPair parses a single "field=kw1,kw2,..." token.
func parseKeywordPair(pair string) (*filter.KeywordFilter, error) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("keyword filter: expected field=kw1,kw2 format, got %q", pair)
	}
	field := parts[0]
	if field == "" {
		return nil, fmt.Errorf("keyword filter: field must not be empty in %q", pair)
	}
	rawKWs := strings.Split(parts[1], ",")
	keywords := make([]string, 0, len(rawKWs))
	for _, k := range rawKWs {
		trimmed := strings.TrimSpace(k)
		if trimmed == "" {
			return nil, fmt.Errorf("keyword filter: empty keyword in %q", pair)
		}
		keywords = append(keywords, trimmed)
	}
	return filter.NewKeywordFilter(field, keywords)
}
