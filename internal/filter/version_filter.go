package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches simple Major.Minor.Patch version strings.
var versionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

type versionFilter struct {
	field string
	min   [3]int
	max   [3]int
	hasMin bool
	hasMax bool
}

// NewVersionFilter returns a Filter that matches entries whose named field
// contains a version string ("Major.Minor.Patch") within [min, max].
// Either bound may be empty to leave it unbounded.
func NewVersionFilter(field, min, max string) (Filter, error) {
	if field == "" {
		return nil, fmt.Errorf("version filter: field must not be empty")
	}
	if min == "" && max == "" {
		return nil, fmt.Errorf("version filter: at least one of min or max must be set")
	}
	f := &versionFilter{field: field}
	if min != "" {
		v, err := parseVersion(min)
		if err != nil {
			return nil, fmt.Errorf("version filter: invalid min %q: %w", min, err)
		}
		f.min = v
		f.hasMin = true
	}
	if max != "" {
		v, err := parseVersion(max)
		if err != nil {
			return nil, fmt.Errorf("version filter: invalid max %q: %w", max, err)
		}
		f.max = v
		f.hasMax = true
	}
	if f.hasMin && f.hasMax && versionCmp(f.min, f.max) > 0 {
		return nil, fmt.Errorf("version filter: min %q is greater than max %q", min, max)
	}
	return f, nil
}

func (f *versionFilter) Match(entry map[string]interface{}) bool {
	val, ok := entry[f.field]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		return false
	}
	v, err := parseVersion(s)
	if err != nil {
		return false
	}
	if f.hasMin && versionCmp(v, f.min) < 0 {
		return false
	}
	if f.hasMax && versionCmp(v, f.max) > 0 {
		return false
	}
	return true
}

func parseVersion(s string) ([3]int, error) {
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("expected Major.Minor.Patch, got %q", s)
	}
	var v [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return [3]int{}, fmt.Errorf("invalid version component %q", p)
		}
		v[i] = n
	}
	return v, nil
}

func versionCmp(a, b [3]int) int {
	for i := 0; i < 3; i++ {
		if a[i] != b[i] {
			if a[i] < b[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}
