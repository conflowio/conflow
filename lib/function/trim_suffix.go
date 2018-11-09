package function

import (
	"strings"
)

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
//go:generate basil generate TrimSuffix
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}
