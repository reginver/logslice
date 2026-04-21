package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Format represents the detected log format.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatText
)

// String returns a human-readable name for the Format.
func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatText:
		return "text"
	default:
		return "unknown"
	}
}

// DetectFormat peeks at the first non-empty line of r to determine whether
// the input looks like JSON or plain text. It returns the detected Format and
// a new reader that includes the consumed bytes so callers can still read the
// full stream.
func DetectFormat(r io.Reader) (Format, io.Reader, error) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)

	scanner := bufio.NewScanner(tee)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// Drain the rest into buf so the combined reader is complete.
		_, _ = io.Copy(buf, r)
		combined := io.MultiReader(bytes.NewReader(buf.Bytes()), strings.NewReader(""))
		if strings.HasPrefix(line, "{") {
			return FormatJSON, combined, nil
		}
		return FormatText, combined, nil
	}
	if err := scanner.Err(); err != nil {
		return FormatUnknown, bytes.NewReader(buf.Bytes()), err
	}
	return FormatUnknown, bytes.NewReader(buf.Bytes()), nil
}
