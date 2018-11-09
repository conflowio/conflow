package variable

import (
	"fmt"
	"strconv"
)

// Number is an integer or float number
type Number struct {
	value interface{}
}

// NewNumber creates a new number
func NewNumber(value interface{}) *Number {
	if unionType, ok := value.(Union); ok {
		value = unionType.Value()
	}
	return &Number{
		value: value,
	}
}

// Value returns with the contained value
func (n *Number) Value() interface{} {
	return n.value
}

// Type returns the type of the value
func (n *Number) Type() string {
	switch n.value.(type) {
	case int64:
		return TypeInteger
	case float64:
		return TypeFloat
	default:
		panic(fmt.Sprintf("unexpected type: %T", n.value))
	}
}

// Float64 returns the float64 value of the variable if it's the correct type
func (n *Number) Float64() (float64, bool) {
	v, ok := n.value.(float64)
	return v, ok
}

// Int64 returns the int64 value of the variable if it's the correct type
func (n *Number) Int64() (int64, bool) {
	v, ok := n.value.(int64)
	return v, ok
}

// String returns with the string representation of the number
func (n *Number) String() string {
	switch v := n.value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	default:
		return ""
	}
}

// IsNumberType returns true if the given value is a number type
func IsNumberType(val interface{}) bool {
	switch val.(type) {
	case int64, float64:
		return true
	default:
		return false
	}
}
