package filter

import (
	"fmt"
	"net"

	"github.com/your/logslice/internal/parser"
)

// IPFilter matches log entries where a field value falls within a CIDR range.
type IPFilter struct {
	field string
	net   *net.IPNet
}

// NewIPFilter creates a filter that matches entries where the given field
// contains an IP address within the specified CIDR block (e.g. "192.168.1.0/24").
func NewIPFilter(field, cidr string) (*IPFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("ip_filter: field name must not be empty")
	}
	if cidr == "" {
		return nil, fmt.Errorf("ip_filter: cidr must not be empty")
	}
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("ip_filter: invalid CIDR %q: %w", cidr, err)
	}
	return &IPFilter{field: field, net: ipNet}, nil
}

// Match returns true if the entry's field is an IP address contained in the CIDR.
func (f *IPFilter) Match(e parser.LogEntry) bool {
	v, ok := e.Fields[f.field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	return f.net.Contains(ip)
}
