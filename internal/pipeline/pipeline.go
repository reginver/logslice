package pipeline

import (
	"io"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

// Config holds the configuration for a pipeline run.
type Config struct {
	From    string
	To      string
	Fields  map[string]string
}

// Result holds the output of a pipeline run.
type Result struct {
	Entries []parser.LogEntry
	Skipped int
}

// Run reads log entries from r, applies time and field filters, and returns
// the matching entries along with a count of skipped lines.
func Run(r io.Reader, cfg Config) (*Result, error) {
	entries, skipped, err := parser.ParseJSON(r)
	if err != nil {
		return nil, err
	}

	tf, err := filter.NewTimeFilter(cfg.From, cfg.To)
	if err != nil {
		return nil, err
	}

	ff, err := filter.NewFieldFilter(cfg.Fields)
	if err != nil {
		return nil, err
	}

	var matched []parser.LogEntry
	for _, e := range entries {
		if !tf.Match(e) {
			continue
		}
		if !ff.Match(e) {
			continue
		}
		matched = append(matched, e)
	}

	return &Result{
		Entries: matched,
		Skipped: skipped,
	}, nil
}
