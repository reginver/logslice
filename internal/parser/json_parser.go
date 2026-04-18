package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogEntry represents a single parsed log line.
type LogEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// TimeField is the JSON key used to extract the timestamp.
const TimeField = "time"

// ParseJSON reads newline-delimited JSON from r and returns a slice of LogEntry.
// Lines that cannot be parsed are skipped with a warning written to errOut.
func ParseJSON(r io.Reader, errOut io.Writer) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			fmt.Fprintf(errOut, "warn: line %d: skipping invalid JSON: %v\n", lineNum, err)
			continue
		}

		entry := LogEntry{Fields: fields, Raw: line}

		if ts, ok := fields[TimeField]; ok {
			switch v := ts.(type) {
			case string:
				parsed, err := time.Parse(time.RFC3339Nano, v)
				if err != nil {
					parsed, err = time.Parse(time.RFC3339, v)
				}
				if err == nil {
					entry.Timestamp = parsed
				}
			case float64:
				entry.Timestamp = time.Unix(int64(v), 0).UTC()
			}
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return entries, nil
}
