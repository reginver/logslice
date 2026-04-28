package cli

import (
	"testing"
)

func TestParseGeoFilter_Valid(t *testing.T) {
	f, err := parseGeoFilter("lat=latitude,lon=longitude,clat=51.5,clon=-0.1,radius=20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestParseGeoFilter_Empty(t *testing.T) {
	f, err := parseGeoFilter("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Fatal("expected nil filter for empty input")
	}
}

func TestParseGeoFilter_MissingLatField(t *testing.T) {
	_, err := parseGeoFilter("lon=longitude,clat=51.5,clon=-0.1,radius=20")
	if err == nil {
		t.Fatal("expected error for missing lat field")
	}
}

func TestParseGeoFilter_MissingLonField(t *testing.T) {
	_, err := parseGeoFilter("lat=latitude,clat=51.5,clon=-0.1,radius=20")
	if err == nil {
		t.Fatal("expected error for missing lon field")
	}
}

func TestParseGeoFilter_MissingCLat(t *testing.T) {
	_, err := parseGeoFilter("lat=latitude,lon=longitude,clon=-0.1,radius=20")
	if err == nil {
		t.Fatal("expected error for missing clat")
	}
}

func TestParseGeoFilter_InvalidRadius(t *testing.T) {
	_, err := parseGeoFilter("lat=latitude,lon=longitude,clat=51.5,clon=-0.1,radius=notanumber")
	if err == nil {
		t.Fatal("expected error for non-numeric radius")
	}
}

func TestParseGeoFilter_ZeroRadius(t *testing.T) {
	_, err := parseGeoFilter("lat=latitude,lon=longitude,clat=51.5,clon=-0.1,radius=0")
	if err == nil {
		t.Fatal("expected error for zero radius")
	}
}

func TestSplitKV_Basic(t *testing.T) {
	m := splitKV("a=1,b=2,c=three")
	if m["a"] != "1" || m["b"] != "2" || m["c"] != "three" {
		t.Errorf("unexpected map: %v", m)
	}
}
