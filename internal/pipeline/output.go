package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/user/logslice/internal/parser"
)

// Format controls how results are written.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Write serialises entries to w in the requested format.
func Write(w io.Writer, entries []parser.LogEntry, fmt_ Format) error {
	switch fmt_ {
	case FormatJSON:
		return writeJSON(w, entries)
	case FormatText:
		return writeText(w, entries)
	default:
		return fmt.Errorf("unknown format %q", fmt_)
	}
}

func writeJSON(w io.Writer, entries []parser.LogEntry) error {
	enc := json.NewEncoder(w)
	for _, e := range entries {
		row := make(map[string]interface{}, len(e.Fields)+1)
		for k, v := range e.Fields {
			row[k] = v
		}
		row["time"] = e.Timestamp.Format("2006-01-02T15:04:05Z07:00")
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	return nil
}

func writeText(w io.Writer, entries []parser.LogEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIME\tFIELDS")
	for _, e := range entries {
		fields, _ := json.Marshal(e.Fields)
		fmt.Fprintf(tw, "%s\t%s\n", e.Timestamp.Format("2006-01-02T15:04:05Z"), fields)
	}
	return tw.Flush()
}
