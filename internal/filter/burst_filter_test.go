package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func burstEntry(field, value string, ts time.Time) parser.Entry {
	return parser.Entry{
		Fields:    map[string]interface{}{field: value},
		Timestamp: &ts,
	}
}

func TestNewBurstFilter_Valid(t *testing.T) {
	_, err := filter.NewBurstFilter("ip", time.Minute, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewBurstFilter_EmptyField(t *testing.T) {
	_, err := filter.NewBurstFilter("", time.Minute, 3)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewBurstFilter_ZeroWindow(t *testing.T) {
	_, err := filter.NewBurstFilter("ip", 0, 3)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewBurstFilter_MinCountTooLow(t *testing.T) {
	_, err := filter.NewBurstFilter("ip", time.Minute, 1)
	if err == nil {
		t.Fatal("expected error for minCount < 2")
	}
}

func TestBurstFilter_NotEnoughOccurrences(t *testing.T) {
	f, _ := filter.NewBurstFilter("ip", time.Minute, 3)
	base := time.Unix(1000, 0)
	e1 := burstEntry("ip", "10.0.0.1", base)
	e2 := burstEntry("ip", "10.0.0.1", base.Add(10*time.Second))
	if f.Match(e1) {
		t.Error("first occurrence should not match")
	}
	if f.Match(e2) {
		t.Error("second occurrence should not match (need 3)")
	}
}

func TestBurstFilter_ReachesThreshold(t *testing.T) {
	f, _ := filter.NewBurstFilter("ip", time.Minute, 3)
	base := time.Unix(2000, 0)
	for i := 0; i < 2; i++ {
		f.Match(burstEntry("ip", "1.2.3.4", base.Add(time.Duration(i)*10*time.Second)))
	}
	third := burstEntry("ip", "1.2.3.4", base.Add(20*time.Second))
	if !f.Match(third) {
		t.Error("third occurrence within window should match")
	}
}

func TestBurstFilter_WindowEviction(t *testing.T) {
	f, _ := filter.NewBurstFilter("ip", 30*time.Second, 3)
	base := time.Unix(3000, 0)
	// two entries inside window
	f.Match(burstEntry("ip", "5.5.5.5", base))
	f.Match(burstEntry("ip", "5.5.5.5", base.Add(5*time.Second)))
	// third entry is outside the window relative to next entry
	newBase := base.Add(60 * time.Second)
	// first entry at newBase — old ones evicted, count resets to 1
	if f.Match(burstEntry("ip", "5.5.5.5", newBase)) {
		t.Error("should not match after window eviction resets count")
	}
}

func TestBurstFilter_MissingField(t *testing.T) {
	f, _ := filter.NewBurstFilter("ip", time.Minute, 2)
	e := parser.Entry{Fields: map[string]interface{}{"other": "val"}, Timestamp: nil}
	if f.Match(e) {
		t.Error("entry without field should not match")
	}
}
