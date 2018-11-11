package function

import (
	"fmt"
	"strconv"

	"github.com/opsidian/basil/variable"
)

// String converts the given value to a string
//go:generate basil generate
func String(value *variable.Basic) string {
	switch v := value.Value().(type) {
	case bool:
		return strconv.FormatBool(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Sprintf("unexpected value type: %T", value.Value()))
	}
}
