package function

import (
	"strings"
)

// Upper returns a copy of the string s with all Unicode letters mapped to their upper case.
//go:generate basil generate
func Upper(s string) string {
	return strings.ToUpper(s)
}
