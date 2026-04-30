package filter

import "testing"

func TestHostnameFilter_WithComposite(t *testing.T) {
	hf, err := NewHostnameFilter("host", []string{"web-*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lf, err := NewLevelFilter("level", "warn", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	comp := NewCompositeFilter(hf, lf)

	pass := map[string]interface{}{"host": "web-01", "level": "error"}
	if !comp.Match(pass) {
		t.Error("expected composite match")
	}

	failHost := map[string]interface{}{"host": "db-01", "level": "error"}
	if comp.Match(failHost) {
		t.Error("expected composite miss on wrong host")
	}

	failLevel := map[string]interface{}{"host": "web-02", "level": "debug"}
	if comp.Match(failLevel) {
		t.Error("expected composite miss on wrong level")
	}
}

func TestHostnameFilter_Negated(t *testing.T) {
	hf, err := NewHostnameFilter("host", []string{"web-*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nf := NewNotFilter(hf)

	if nf.Match(hostnameEntry("host", "web-01")) {
		t.Error("expected negated filter to reject web-01")
	}
	if !nf.Match(hostnameEntry("host", "db-01")) {
		t.Error("expected negated filter to pass db-01")
	}
}
