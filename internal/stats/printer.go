package stats

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

const timeLayout = "2006-01-02T15:04:05Z"

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Summary) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "--- Log Slice Summary ---")
	fmt.Fprintf(tw, "Total lines:\t%d\n", s.TotalLines)
	fmt.Fprintf(tw, "Matched lines:\t%d\n", s.MatchedLines)
	fmt.Fprintf(tw, "Skipped lines:\t%d\n", s.SkippedLines)

	if s.EarliestTime != nil {
		fmt.Fprintf(tw, "Earliest match:\t%s\n", s.EarliestTime.UTC().Format(timeLayout))
	}
	if s.LatestTime != nil {
		fmt.Fprintf(tw, "Latest match:\t%s\n", s.LatestTime.UTC().Format(timeLayout))
	}

	if len(s.FieldCounts) > 0 {
		fmt.Fprintln(tw, "\nField value counts:")
		keys := make([]string, 0, len(s.FieldCounts))
		for k := range s.FieldCounts {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(tw, "  %s:\t%d\n", k, s.FieldCounts[k])
		}
	}

	tw.Flush()
}
