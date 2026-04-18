// Package stats provides pipeline run statistics collection and reporting.
//
// Use NewCollector to create a Collector, record events during processing,
// then call Summary to retrieve the aggregated results. Use Print to render
// a human-readable summary to any io.Writer.
package stats
