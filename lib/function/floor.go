package function

import "math"

// Floor returns the greatest integer value less than or equal to x.
//go:generate basil generate Floor
func Floor(value float64) int64 {
	return int64(math.Floor(value))
}
