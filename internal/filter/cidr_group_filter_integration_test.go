package filter

import (
	"testing"
)

// TestCIDRGroupFilter_WithComposite verifies that CIDRGroupFilter composes
// correctly with CompositeFilter (all-match semantics).
func TestCIDRGroupFilter_WithComposite(t *testing.T) {
	cg, err := NewCIDRGroupFilter("ip", []string{"10.0.0.0/8"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	pf, err := NewPrefixFilter("env", "prod")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	comp := NewCompositeFilter(cg, pf)

	hit := map[string]interface{}{"ip": "10.2.3.4", "env": "production"}
	if !comp.Match(hit) {
		t.Error("expected composite match")
	}

	miss := map[string]interface{}{"ip": "172.16.0.1", "env": "production"}
	if comp.Match(miss) {
		t.Error("expected composite miss when IP is outside range")
	}
}

// TestCIDRGroupFilter_Negated verifies that wrapping CIDRGroupFilter in
// NotFilter inverts the match result.
func TestCIDRGroupFilter_Negated(t *testing.T) {
	cg, err := NewCIDRGroupFilter("ip", []string{"10.0.0.0/8"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	not, err := NewNotFilter(cg)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	internalIP := map[string]interface{}{"ip": "10.0.0.1"}
	if not.Match(internalIP) {
		t.Error("expected NotFilter to reject internal IP")
	}

	externalIP := map[string]interface{}{"ip": "8.8.8.8"}
	if !not.Match(externalIP) {
		t.Error("expected NotFilter to pass external IP")
	}
}
