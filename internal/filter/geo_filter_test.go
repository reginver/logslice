package filter

import (
	"testing"
)

func geoEntry(lat, lon float64) map[string]interface{} {
	return map[string]interface{}{"lat": lat, "lon": lon}
}

func TestNewGeoFilter_Valid(t *testing.T) {
	f, err := NewGeoFilter("lat", "lon", 51.5, -0.1, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewGeoFilter_EmptyLatField(t *testing.T) {
	_, err := NewGeoFilter("", "lon", 0, 0, 10)
	if err == nil {
		t.Fatal("expected error for empty latField")
	}
}

func TestNewGeoFilter_EmptyLonField(t *testing.T) {
	_, err := NewGeoFilter("lat", "", 0, 0, 10)
	if err == nil {
		t.Fatal("expected error for empty lonField")
	}
}

func TestNewGeoFilter_ZeroRadius(t *testing.T) {
	_, err := NewGeoFilter("lat", "lon", 0, 0, 0)
	if err == nil {
		t.Fatal("expected error for zero radius")
	}
}

func TestNewGeoFilter_NegativeRadius(t *testing.T) {
	_, err := NewGeoFilter("lat", "lon", 0, 0, -5)
	if err == nil {
		t.Fatal("expected error for negative radius")
	}
}

func TestGeoFilter_Match_WithinRadius(t *testing.T) {
	// London ~51.5074, -0.1278; radius 10 km
	f, _ := NewGeoFilter("lat", "lon", 51.5074, -0.1278, 10)
	// Covent Garden: 51.5117, -0.1240 (~0.5 km away)
	if !f.Match(geoEntry(51.5117, -0.1240)) {
		t.Error("expected match for nearby point")
	}
}

func TestGeoFilter_Match_OutsideRadius(t *testing.T) {
	f, _ := NewGeoFilter("lat", "lon", 51.5074, -0.1278, 10)
	// Paris: 48.8566, 2.3522 (~340 km away)
	if f.Match(geoEntry(48.8566, 2.3522)) {
		t.Error("expected no match for distant point")
	}
}

func TestGeoFilter_Match_MissingField(t *testing.T) {
	f, _ := NewGeoFilter("lat", "lon", 51.5, -0.1, 50)
	if f.Match(map[string]interface{}{"lat": 51.5}) {
		t.Error("expected no match when lon field is absent")
	}
}

func TestGeoFilter_Match_NonNumericField(t *testing.T) {
	f, _ := NewGeoFilter("lat", "lon", 51.5, -0.1, 50)
	if f.Match(map[string]interface{}{"lat": "north", "lon": -0.1}) {
		t.Error("expected no match for non-numeric lat")
	}
}

func TestGeoFilter_Match_ExactCentre(t *testing.T) {
	f, _ := NewGeoFilter("lat", "lon", 48.8566, 2.3522, 1)
	if !f.Match(geoEntry(48.8566, 2.3522)) {
		t.Error("expected match at exact centre")
	}
}
