package filter

import (
	"fmt"
	"math"
)

// GeoFilter matches log entries whose lat/lon fields fall within a given radius
// (in kilometres) of a centre point.
type GeoFilter struct {
	latField string
	lonField string
	centreLat float64
	centreLon float64
	radiusKm  float64
}

const earthRadiusKm = 6371.0

// NewGeoFilter returns a GeoFilter that passes entries where the great-circle
// distance from (centreLat, centreLon) is ≤ radiusKm.
func NewGeoFilter(latField, lonField string, centreLat, centreLon, radiusKm float64) (*GeoFilter, error) {
	if latField == "" {
		return nil, fmt.Errorf("geo filter: latField must not be empty")
	}
	if lonField == "" {
		return nil, fmt.Errorf("geo filter: lonField must not be empty")
	}
	if radiusKm <= 0 {
		return nil, fmt.Errorf("geo filter: radiusKm must be positive, got %g", radiusKm)
	}
	return &GeoFilter{
		latField:  latField,
		lonField:  lonField,
		centreLat: centreLat,
		centreLon: centreLon,
		radiusKm:  radiusKm,
	}, nil
}

// Match returns true when the entry's lat/lon are within the configured radius.
func (f *GeoFilter) Match(entry map[string]interface{}) bool {
	lat, ok := toFloat64(entry[f.latField])
	if !ok {
		return false
	}
	lon, ok := toFloat64(entry[f.lonField])
	if !ok {
		return false
	}
	return haversine(f.centreLat, f.centreLon, lat, lon) <= f.radiusKm
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func toRad(deg float64) float64 { return deg * math.Pi / 180 }

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}
