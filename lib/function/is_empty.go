package function

import (
	"github.com/opsidian/basil/variable"
)

// IsEmpty returns true if the given value has an empty value
//go:generate basil generate
func IsEmpty(value interface{}) bool {
	return variable.IsEmpty(value)
}
