package function

import "math"

// Ceil returns the least integer value greater than or equal to x.
//go:generate basil generate Ceil
func Ceil(value float64) int64 {
	return int64(math.Ceil(value))
}
