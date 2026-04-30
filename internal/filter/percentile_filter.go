package filter

import (
	"encoding/json"
	"fmt"
	"sort"
)

// PercentileFilter keeps only entries whose numeric field value falls within
// the [minPct, maxPct] percentile band, computed over all seen values.
type PercentileFilter struct {
	field  string
	minPct float64
	maxPct float64
	values []float64
}

// NewPercentileFilter returns a filter that passes entries whose field value
// lies within [minPct, maxPct] (0–100) of the observed distribution.
func NewPercentileFilter(field string, minPct, maxPct float64) (*PercentileFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("percentile filter: field must not be empty")
	}
	if minPct < 0 || maxPct > 100 {
		return nil, fmt.Errorf("percentile filter: percentiles must be in [0, 100]")
	}
	if minPct >= maxPct {
		return nil, fmt.Errorf("percentile filter: minPct must be less than maxPct")
	}
	return &PercentileFilter{
		field:  field,
		minPct: minPct,
		maxPct: maxPct,
	}, nil
}

// Match records the entry's field value and returns true if it falls within
// the configured percentile range based on all values seen so far.
func (f *PercentileFilter) Match(entry map[string]interface{}) bool {
	val, ok := entry[f.field]
	if !ok {
		return false
	}
	v, err := pctToFloat64(val)
	if err != nil {
		return false
	}
	f.values = append(f.values, v)

	n := len(f.values)
	sorted := make([]float64, n)
	copy(sorted, f.values)
	sort.Float64s(sorted)

	rank := sort.SearchFloat64s(sorted, v)
	pct := float64(rank) / float64(n) * 100.0
	return pct >= f.minPct && pct <= f.maxPct
}

func pctToFloat64(v interface{}) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case json.Number:
		return n.Float64()
	}
	return 0, fmt.Errorf("not a number: %T", v)
}
