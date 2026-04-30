package filter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ChecksumFilter matches log entries where a field value, when hashed with the
// specified algorithm, produces a digest that starts with the given prefix.
type ChecksumFilter struct {
	field  string
	algo   string
	prefix string
}

// NewChecksumFilter returns a filter that hashes entry[field] with algo and
// checks whether the resulting hex digest has the given prefix.
// algo must be one of "md5", "sha1", or "sha256".
func NewChecksumFilter(field, algo, prefix string) (*ChecksumFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("checksum filter: field must not be empty")
	}
	algo = strings.ToLower(algo)
	switch algo {
	case "md5", "sha1", "sha256":
	default:
		return nil, fmt.Errorf("checksum filter: unknown algorithm %q (want md5, sha1, sha256)", algo)
	}
	if prefix == "" {
		return nil, fmt.Errorf("checksum filter: prefix must not be empty")
	}
	return &ChecksumFilter{field: field, algo: algo, prefix: prefix}, nil
}

// Match returns true when the hashed field value starts with the configured prefix.
func (f *ChecksumFilter) Match(entry map[string]interface{}) bool {
	v, ok := entry[f.field]
	if !ok {
		return false
	}
	s := fmt.Sprintf("%v", v)
	digest := f.hash(s)
	return strings.HasPrefix(digest, strings.ToLower(f.prefix))
}

func (f *ChecksumFilter) hash(s string) string {
	b := []byte(s)
	switch f.algo {
	case "sha1":
		h := sha1.Sum(b)
		return hex.EncodeToString(h[:])
	case "sha256":
		h := sha256.Sum256(b)
		return hex.EncodeToString(h[:])
	default: // md5
		h := md5.Sum(b)
		return hex.EncodeToString(h[:])
	}
}
