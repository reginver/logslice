package filter

import (
	"crypto/md5"
	"fmt"
	"errors"
)

// HashFilter matches log entries where a field's MD5 hash starts with a given prefix.
// This is useful for consistent sampling or routing based on field identity.
type HashFilter struct {
	field  string
	prefix string
}

// NewHashFilter returns a HashFilter that matches entries where the MD5 hex
// digest of the named field's string value starts with prefix.
func NewHashFilter(field, prefix string) (*HashFilter, error) {
	if field == "" {
		return nil, errors.New("hash filter: field name must not be empty")
	}
	if prefix == "" {
		return nil, errors.New("hash filter: prefix must not be empty")
	}
	if len(prefix) > 32 {
		return nil, errors.New("hash filter: prefix length exceeds MD5 hex length (32)")
	}
	return &HashFilter{field: field, prefix: prefix}, nil
}

// Match returns true when the MD5 hex digest of the entry's field value
// starts with the configured prefix.
func (f *HashFilter) Match(entry map[string]interface{}) bool {
	v, ok := entry[f.field]
	if !ok {
		return false
	}
	str := fmt.Sprintf("%v", v)
	sum := md5.Sum([]byte(str))
	hex := fmt.Sprintf("%x", sum)
	return len(hex) >= len(f.prefix) && hex[:len(f.prefix)] == f.prefix
}
