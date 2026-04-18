package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/pipeline"
	"github.com/user/logslice/internal/stats"
)

// Run is the entry point for the CLI.
func Run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	in, err := openInput(cfg.Input)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	if f, ok := in.(*os.File); ok && f != os.Stdin {
		defer f.Close()
	}

	out, err := openOutput(cfg.Output)
	if err != nil {
		return fmt.Errorf("opening output: %w", err)
	}
	if f, ok := out.(*os.File); ok && f != os.Stdout {
		defer f.Close()
	}

	var filters []pipeline.Filter

	if cfg.TimeStart != "" || cfg.TimeEnd != "" {
		tf, err := filter.NewTimeFilter(cfg.TimeStart, cfg.TimeEnd)
		if err != nil {
			return fmt.Errorf("time filter: %w", err)
		}
		filters = append(filters, tf)
	}

	for k, v := range cfg.Fields {
		ff, err := filter.NewFieldFilter(k, v)
		if err != nil {
			return fmt.Errorf("field filter: %w", err)
		}
		filters = append(filters, ff)
	}

	collector := stats.NewCollector()

	if err := pipeline.Run(in, out, cfg.Format, filters, collector); err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}

	if cfg.ShowStats {
		stats.Print(os.Stderr, collector)
	}

	return nil
}

func openInput(path string) (io.Reader, error) {
	if path == "" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func openOutput(path string) (io.Writer, error) {
	if path == "" {
		return os.Stdout, nil
	}
	return os.Create(path)
}
