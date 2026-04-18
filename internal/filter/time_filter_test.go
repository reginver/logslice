package filter

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNewTimeFilter_Valid(t *testing.T) {
	tf, err := NewTimeFilter("2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tf.Start == nil || tf.End == nil {
		t.Fatal("expected both start and end to be set")
	}
}

func TestNewTimeFilter_EmptyBounds(t *testing.T) {
	tf, err := NewTimeFilter("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tf.Start != nil || tf.End != nil {
		t.Fatal("expected nil bounds for empty input")
	}
}

func TestNewTimeFilter_EndBeforeStart(t *testing.T) {
	_, err := NewTimeFilter("2024-01-02T00:00:00Z", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error when end is before start")
	}
}

func TestNewTimeFilter_InvalidFormat(t *testing.T) {
	_, err := NewTimeFilter("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid start format")
	}
}

func TestTimeFilter_Match(t *testing.T) {
	tf, _ := NewTimeFilter("2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")

	cases := []struct {
		ts      string
		expect  bool
	}{
		{"2024-01-01T11:00:00Z", true},
		{"2024-01-01T10:00:00Z", true},
		{"2024-01-01T12:00:00Z", true},
		{"2024-01-01T09:59:59Z", false},
		{"2024-01-01T12:00:01Z", false},
	}

	for _, c := range cases {
		got := tf.Match(mustTime(c.ts))
		if got != c.expect {
			t.Errorf("Match(%s) = %v, want %v", c.ts, got, c.expect)
		}
	}
}

func TestTimeFilter_Match_NoBounds(t *testing.T) {
	tf, _ := NewTimeFilter("", "")
	if !tf.Match(time.Now()) {
		t.Error("unbounded filter should match any time")
	}
}
