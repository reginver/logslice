package filter

import (
	"testing"
)

// TestGeoFilter_WithComposite verifies GeoFilter composes correctly inside a
// CompositeFilter (all-match) alongside a FieldFilter.
func TestGeoFilter_WithComposite(t *testing.T) {
	geo, err := NewGeoFilter("lat", "lon", 40.7128, -74.0060, 50) // NYC, 50 km
	if err != nil {
		t.Fatalf("NewGeoFilter: %v", err)
	}
	field, err := NewFieldFilter("env", "production")
	if err != nil {
		t.Fatalf("NewFieldFilter: %v", err)
	}
	comp := NewCompositeFilter(geo, field)

	// Newark, NJ: ~16 km from NYC — should pass both
	match := map[string]interface{}{"lat": 40.7357, "lon": -74.1724, "env": "production"}
	if !comp.Match(match) {
		t.Error("expected composite match for nearby production entry")
	}

	// Wrong env
	wrongEnv := map[string]interface{}{"lat": 40.7357, "lon": -74.1724, "env": "staging"}
	if comp.Match(wrongEnv) {
		t.Error("expected no match when env differs")
	}

	// Too far (Boston: ~306 km)
	farAway := map[string]interface{}{"lat": 42.3601, "lon": -71.0589, "env": "production"}
	if comp.Match(farAway) {
		t.Error("expected no match for distant entry")
	}
}

// TestGeoFilter_Negated verifies that NotFilter correctly inverts GeoFilter.
func TestGeoFilter_Negated(t *testing.T) {
	geo, err := NewGeoFilter("lat", "lon", 40.7128, -74.0060, 50)
	if err != nil {
		t.Fatalf("NewGeoFilter: %v", err)
	}
	not, err := NewNotFilter(geo)
	if err != nil {
		t.Fatalf("NewNotFilter: %v", err)
	}

	nearby := map[string]interface{}{"lat": 40.7357, "lon": -74.1724}
	if not.Match(nearby) {
		t.Error("expected negated filter to reject nearby entry")
	}

	farAway := map[string]interface{}{"lat": 48.8566, "lon": 2.3522}
	if !not.Match(farAway) {
		t.Error("expected negated filter to pass distant entry")
	}
}
