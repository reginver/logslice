package filter

import (
	"fmt"
	"net"
)

// CIDRGroupFilter matches log entries where a named field's IP value
// falls within ANY of the provided CIDR blocks.
type CIDRGroupFilter struct {
	field  string
	nets   []*net.IPNet
}

// NewCIDRGroupFilter returns a filter that passes entries whose field
// value (as a string IP address) is contained in at least one of the
// given CIDR blocks.
func NewCIDRGroupFilter(field string, cidrs []string) (*CIDRGroupFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("cidr_group_filter: field must not be empty")
	}
	if len(cidrs) == 0 {
		return nil, fmt.Errorf("cidr_group_filter: at least one CIDR block required")
	}

	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		if cidr == "" {
			return nil, fmt.Errorf("cidr_group_filter: CIDR block must not be empty")
		}
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("cidr_group_filter: invalid CIDR %q: %w", cidr, err)
		}
		nets = append(nets, ipNet)
	}

	return &CIDRGroupFilter{field: field, nets: nets}, nil
}

// Match returns true if the entry's field value is an IP contained in
// any of the configured CIDR blocks.
func (f *CIDRGroupFilter) Match(entry map[string]interface{}) bool {
	val, ok := entry[f.field]
	if !ok {
		return false
	}
	s, ok := val.(string)
	if !ok {
		return false
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	for _, n := range f.nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
