package function

import (
	"fmt"

	"github.com/opsidian/basil/variable"
)

// Abs returns the absolute value of the given number
//go:generate basil generate
func Abs(value *variable.Number) (*variable.Number, error) {
	switch n := value.Value().(type) {
	case int64:
		if n >= 0 {
			return value, nil
		}
		return variable.NewNumber(-1 * n), nil
	case float64:
		if n >= 0 {
			return value, nil
		}
		return variable.NewNumber(-1 * n), nil
	default:
		panic(fmt.Sprintf("unexpected value type: %T", value.Value()))
	}
}
