package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andrewstuart/logslice/internal/parser"
)

// semver holds a parsed semantic version.
type semver struct {
	major, minor, patch int
}

func parseSemver(s string) (semver, error) {
	s = strings.TrimPrefix(s, "v")
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return semver{}, fmt.Errorf("semver: expected major.minor.patch, got %q", s)
	}
	var sv semver
	var err error
	if sv.major, err = strconv.Atoi(parts[0]); err != nil {
		return semver{}, fmt.Errorf("semver: invalid major %q", parts[0])
	}
	if sv.minor, err = strconv.Atoi(parts[1]); err != nil {
		return semver{}, fmt.Errorf("semver: invalid minor %q", parts[1])
	}
	patch := strings.SplitN(parts[2], "-", 2)[0] // strip pre-release
	if sv.patch, err = strconv.Atoi(patch); err != nil {
		return semver{}, fmt.Errorf("semver: invalid patch %q", patch)
	}
	return sv, nil
}

func (a semver) less(b semver) bool {
	if a.major != b.major {
		return a.major < b.major
	}
	if a.minor != b.minor {
		return a.minor < b.minor
	}
	return a.patch < b.patch
}

// SemverFilter matches log entries whose named field falls within [min, max].
// Either bound may be empty to indicate unbounded.
type SemverFilter struct {
	field    string
	hasMin   bool
	hasMax   bool
	min, max semver
}

// NewSemverFilter constructs a SemverFilter. minVer and maxVer may be empty.
func NewSemverFilter(field, minVer, maxVer string) (*SemverFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("semver filter: field must not be empty")
	}
	if minVer == "" && maxVer == "" {
		return nil, fmt.Errorf("semver filter: at least one of min or max must be set")
	}
	f := &SemverFilter{field: field}
	if minVer != "" {
		sv, err := parseSemver(minVer)
		if err != nil {
			return nil, err
		}
		f.min = sv
		f.hasMin = true
	}
	if maxVer != "" {
		sv, err := parseSemver(maxVer)
		if err != nil {
			return nil, err
		}
		f.max = sv
		f.hasMax = true
	}
	if f.hasMin && f.hasMax && f.max.less(f.min) {
		return nil, fmt.Errorf("semver filter: max %s is less than min %s", maxVer, minVer)
	}
	return f, nil
}

// Match returns true when the entry's field value is a semver within [min, max].
func (f *SemverFilter) Match(e parser.LogEntry) bool {
	raw, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := raw.(string)
	if !ok {
		return false
	}
	sv, err := parseSemver(s)
	if err != nil {
		return false
	}
	if f.hasMin && sv.less(f.min) {
		return false
	}
	if f.hasMax && f.max.less(sv) {
		return false
	}
	return true
}
