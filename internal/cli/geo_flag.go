package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// parseGeoFilter parses a geo flag value of the form:
//   lat=<field>,lon=<field>,clat=<float>,clon=<float>,radius=<km>
func parseGeoFilter(raw string) (*filter.GeoFilter, error) {
	if raw == "" {
		return nil, nil
	}
	parts := splitKV(raw)
	latField := parts["lat"]
	lonField := parts["lon"]
	if latField == "" {
		return nil, fmt.Errorf("geo: missing lat=<field>")
	}
	if lonField == "" {
		return nil, fmt.Errorf("geo: missing lon=<field>")
	}
	clat, err := parseFloatKey(parts, "clat")
	if err != nil {
		return nil, fmt.Errorf("geo: %w", err)
	}
	clon, err := parseFloatKey(parts, "clon")
	if err != nil {
		return nil, fmt.Errorf("geo: %w", err)
	}
	radius, err := parseFloatKey(parts, "radius")
	if err != nil {
		return nil, fmt.Errorf("geo: %w", err)
	}
	return filter.NewGeoFilter(latField, lonField, clat, clon, radius)
}

// splitKV splits "k1=v1,k2=v2" into a map.
func splitKV(s string) map[string]string {
	m := make(map[string]string)
	for _, tok := range strings.Split(s, ",") {
		if idx := strings.IndexByte(tok, '='); idx >= 0 {
			m[tok[:idx]] = tok[idx+1:]
		}
	}
	return m
}

func parseFloatKey(m map[string]string, key string) (float64, error) {
	v, ok := m[key]
	if !ok || v == "" {
		return 0, fmt.Errorf("missing %s=<float>", key)
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", key, v, err)
	}
	return f, nil
}
