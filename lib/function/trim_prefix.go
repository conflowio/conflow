package function

import (
	"strings"
)

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
//go:generate basil generate TrimPrefix
func TrimPrefix(s string, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
