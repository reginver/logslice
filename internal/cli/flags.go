package cli

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds all parsed CLI flags.
type Config struct {
	Input      string
	Output     string
	Format     string
	Start      string
	End        string
	Fields     map[string]string // field=value exact matches
	Regexes    map[string]string // field=pattern regex matches
	Stats      bool
	OutputFmt  string
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	input := fs.String("input", "", "input log file (default: stdin)")
	output := fs.String("output", "", "output file (default: stdout)")
	format := fs.String("format", "auto", "log format: auto|json|text")
	start := fs.String("start", "", "start time (RFC3339)")
	end := fs.String("end", "", "end time (RFC3339)")
	stats := fs.Bool("stats", false, "print stats after processing")
	outFmt := fs.String("output-format", "json", "output format: json|text")

	var fieldFlags []string
	fs.Func("field", "field=value filter (repeatable)", func(s string) error {
		fieldFlags = append(fieldFlags, s)
		return nil
	})

	var regexFlags []string
	fs.Func("regex", "field=pattern regex filter (repeatable)", func(s string) error {
		regexFlags = append(regexFlags, s)
		return nil
	})

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	fields, err := parsePairs(fieldFlags)
	if err != nil {
		return nil, fmt.Errorf("--field: %w", err)
	}
	regexes, err := parsePairs(regexFlags)
	if err != nil {
		return nil, fmt.Errorf("--regex: %w", err)
	}

	return &Config{
		Input:     *input,
		Output:    *output,
		Format:    *format,
		Start:     *start,
		End:       *end,
		Fields:    fields,
		Regexes:   regexes,
		Stats:     *stats,
		OutputFmt: *outFmt,
	}, nil
}

func parsePairs(raw []string) (map[string]string, error) {
	out := make(map[string]string, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid format %q, expected field=value", r)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}
