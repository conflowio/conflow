package math

import (
	"fmt"
	"math"

	"github.com/opsidian/basil/basil/variable"
)

// Ceil returns the least integer value greater than or equal to x.
//go:generate basil generate
func Ceil(number *variable.Number) int64 {
	switch v := number.Value().(type) {
	case int64:
		return v
	case float64:
		return int64(math.Ceil(v))
	default:
		panic(fmt.Sprintf("unexpected type: %T", number.Value()))
	}
}
