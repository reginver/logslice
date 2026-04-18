package stats

import (
	"time"
)

// Summary holds aggregate statistics collected during a pipeline run.
type Summary struct {
	TotalLines   int
	MatchedLines int
	SkippedLines int
	EarliestTime *time.Time
	LatestTime   *time.Time
	FieldCounts  map[string]int
}

// Collector gathers statistics from log entries.
type Collector struct {
	summary Summary
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{
		summary: Summary{
			FieldCounts: make(map[string]int),
		},
	}
}

// RecordTotal increments the total line counter.
func (c *Collector) RecordTotal() {
	c.summary.TotalLines++
}

// RecordMatched increments the matched line counter and updates time bounds.
func (c *Collector) RecordMatched(ts *time.Time) {
	c.summary.MatchedLines++
	if ts != nil {
		if c.summary.EarliestTime == nil || ts.Before(*c.summary.EarliestTime) {
			t := *ts
			c.summary.EarliestTime = &t
		}
		if c.summary.LatestTime == nil || ts.After(*c.summary.LatestTime) {
			t := *ts
			c.summary.LatestTime = &t
		}
	}
}

// RecordSkipped increments the skipped line counter.
func (c *Collector) RecordSkipped() {
	c.summary.SkippedLines++
}

// RecordField increments the count for a specific field value.
func (c *Collector) RecordField(key string) {
	c.summary.FieldCounts[key]++
}

// Summary returns the collected statistics.
func (c *Collector) Summary() Summary {
	return c.summary
}
