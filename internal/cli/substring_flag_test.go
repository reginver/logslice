package cli

import (
	"testing"
)

func TestParseSubstringPairs_Valid(t *testing.T) {
	specs, err := parseSubstringPairs([]string{"msg:error", "svc:auth:i"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(specs))
	}
	if specs[0].Field != "msg" || specs[0].Substring != "error" || specs[0].CaseFold {
		t.Errorf("unexpected spec[0]: %+v", specs[0])
	}
	if specs[1].Field != "svc" || specs[1].Substring != "auth" || !specs[1].CaseFold {
		t.Errorf("unexpected spec[1]: %+v", specs[1])
	}
}

func TestParseSubstringPairs_Empty(t *testing.T) {
	specs, err := parseSubstringPairs(nil)
	if err != nil || len(specs) != 0 {
		t.Fatalf("expected empty result, got %v %v", specs, err)
	}
}

func TestParseSubstringPair_MissingParts(t *testing.T) {
	_, err := parseSubstringPair("onlyone")
	if err == nil {
		t.Fatal("expected error for missing parts")
	}
}

func TestParseSubstringPair_EmptyField(t *testing.T) {
	_, err := parseSubstringPair(":error")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseSubstringPair_EmptySubstring(t *testing.T) {
	_, err := parseSubstringPair("msg:")
	if err == nil {
		t.Fatal("expected error for empty substring")
	}
}

func TestParseSubstringPair_CaseFoldFlag(t *testing.T) {
	spec, err := parseSubstringPair("msg:warn:I")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !spec.CaseFold {
		t.Error("expected CaseFold=true for :I suffix")
	}
}
