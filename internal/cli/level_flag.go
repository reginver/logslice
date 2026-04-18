package cli

import (
	"fmt"
	"strings"
)

// LevelRange holds the parsed min/max level strings from CLI flags.
type LevelRange struct {
	Field    string
	MinLevel string
	MaxLevel string
}

// parseLevelFlag parses a level range flag value of the form
// "field:minLevel:maxLevel". Either minLevel or maxLevel may be omitted but
// the field and at least one bound must be present.
//
// Examples:
//
//	"level:info:error"  -> field=level, min=info, max=error
//	"level:warn:"       -> field=level, min=warn, max=<open>
//	"level::error"      -> field=level, min=<open>, max=error
func parseLevelFlag(s string) (LevelRange, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return LevelRange{}, fmt.Errorf("level flag: expected field:min:max, got %q", s)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return LevelRange{}, fmt.Errorf("level flag: field name must not be empty")
	}
	min := strings.TrimSpace(parts[1])
	max := strings.TrimSpace(parts[2])
	if min == "" && max == "" {
		return LevelRange{}, fmt.Errorf("level flag: at least one of min or max level must be specified")
	}
	return LevelRange{Field: field, MinLevel: min, MaxLevel: max}, nil
}

// parseLevelFlags parses multiple level flag values.
func parseLevelFlags(vals []string) ([]LevelRange, error) {
	out := make([]LevelRange, 0, len(vals))
	for _, v := range vals {
		lr, err := parseLevelFlag(v)
		if err != nil {
			return nil, err
		}
		out = append(out, lr)
	}
	return out, nil
}

// String returns the canonical string representation of a LevelRange,
// suitable for use as a CLI flag value.
func (lr LevelRange) String() string {
	return fmt.Sprintf("%s:%s:%s", lr.Field, lr.MinLevel, lr.MaxLevel)
}
