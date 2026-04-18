package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

// TextEntry represents a parsed line from a plain-text log.
type TextEntry struct {
	Timestamp *time.Time
	Level     string
	Message   string
	Raw       string
}

// defaultTextPattern matches lines like: 2024-01-02T15:04:05 ERROR some message
var defaultTextPattern = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})?)\s+(\w+)\s+(.+)$`,
)

// ParseText reads lines from r and attempts to parse each as a structured
// text log entry. Lines that do not match the pattern are returned with only
// the Raw field populated.
func ParseText(r io.Reader) ([]TextEntry, error) {
	var entries []TextEntry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		entry := TextEntry{Raw: line}
		if m := defaultTextPattern.FindStringSubmatch(line); m != nil {
			t, err := time.Parse(time.RFC3339, m[1])
			if err == nil {
				entry.Timestamp = &t
			}
			entry.Level = m[2]
			entry.Message = m[3]
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("text parser: %w", err)
	}
	return entries, nil
}
