package cli

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds parsed CLI flags.
type Config struct {
	Input      string
	Output     string
	Format     string
	TimeStart  string
	TimeEnd    string
	Fields     map[string]string
	ShowStats  bool
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	input := fs.String("input", "", "input log file (default: stdin)")
	output := fs.String("output", "", "output file (default: stdout)")
	format := fs.String("format", "json", "output format: json or text")
	timeStart := fs.String("from", "", "start time (RFC3339 or unix timestamp)")
	timeEnd := fs.String("to", "", "end time (RFC3339 or unix timestamp)")
	showStats := fs.Bool("stats", false, "print stats summary to stderr")

	var rawFields []string
	fs.Func("field", "field filter as key=value (repeatable)", func(s string) error {
		rawFields = append(rawFields, s)
		return nil
	})

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	for _, f := range rawFields {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field filter %q: expected key=value", f)
		}
		fields[parts[0]] = parts[1]
	}

	if *format != "json" && *format != "text" {
		return nil, fmt.Errorf("unsupported format %q: must be json or text", *format)
	}

	return &Config{
		Input:     *input,
		Output:    *output,
		Format:    *format,
		TimeStart: *timeStart,
		TimeEnd:   *timeEnd,
		Fields:    fields,
		ShowStats: *showStats,
	}, nil
}
