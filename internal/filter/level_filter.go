package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// levelOrder maps log level names to numeric severity.
var levelOrder = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
}

// LevelFilter passes entries whose level is >= minLevel and <= maxLevel.
type LevelFilter struct {
	field    string
	minLevel int
	maxLevel int
}

// NewLevelFilter creates a LevelFilter. minLevel or maxLevel may be empty to
// leave that bound open.
func NewLevelFilter(field, minLevel, maxLevel string) (*LevelFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("level filter: field name must not be empty")
	}
	min := 0
	max := len(levelOrder) - 1
	if minLevel != "" {
		v, ok := levelOrder[strings.ToLower(minLevel)]
		if !ok {
			return nil, fmt.Errorf("level filter: unknown min level %q", minLevel)
		}
		min = v
	}
	if maxLevel != "" {
		v, ok := levelOrder[strings.ToLower(maxLevel)]
		if !ok {
			return nil, fmt.Errorf("level filter: unknown max level %q", maxLevel)
		}
		max = v
	}
	if min > max {
		return nil, fmt.Errorf("level filter: min level %q is above max level %q", minLevel, maxLevel)
	}
	return &LevelFilter{field: field, minLevel: min, maxLevel: max}, nil
}

// Match returns true when the entry's level field falls within the configured range.
func (f *LevelFilter) Match(e parser.LogEntry) bool {
	raw, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := raw.(string)
	if !ok {
		return false
	}
	v, ok := levelOrder[strings.ToLower(s)]
	if !ok {
		return false
	}
	return v >= f.minLevel && v <= f.maxLevel
}
